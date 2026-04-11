import { useState, useEffect } from 'react'
import {
  Box,
  Typography,
  Stack,
  Paper,
  Card,
  CardContent,
  CardActions,
  Chip,
  CircularProgress,
  Alert,
  Grid,
  Divider,
  Button,
  Snackbar,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  List,
  ListItem,
  Avatar,
} from '@mui/material'
import EventIcon from '@mui/icons-material/Event'
import LocationOnIcon from '@mui/icons-material/LocationOn'
import SportsIcon from '@mui/icons-material/Sports'
import GroupsIcon from '@mui/icons-material/Groups'
import AccessTimeIcon from '@mui/icons-material/AccessTime'
import DeleteIcon from '@mui/icons-material/Delete'
import InboxIcon from '@mui/icons-material/Inbox'
import AddIcon from '@mui/icons-material/Add'
import PersonIcon from '@mui/icons-material/Person'
import { useNavigate } from 'react-router-dom'
import { useMatchOfferContext } from '../../context/MatchOfferContext'
import { fetchAccount, Account } from '../../../auth/infrastructure/adapters/AccountApiAdapter'
import { useMatchRequestContext } from '../../../matchrequest/context/MatchRequestContext'
import CheckIcon from '@mui/icons-material/Check'
import CloseIcon from '@mui/icons-material/Close'
import { MatchOffer } from '../../../../shared/types/matchOffer'
import { MatchRequest } from '../../../matchrequest/domain/ports/MatchRequestRepository'
import { useAuth } from '../../../auth/context/AuthContext'

export function MyOffersPage() {
  const { findAccountMatchOffersUseCase, deleteMatchOfferUseCase } = useMatchOfferContext()
  const { findReceivedMatchRequestsUseCase, acceptMatchRequestUseCase } = useMatchRequestContext()
  const { accountId } = useAuth()
  const navigate = useNavigate()

  const [offers, setOffers] = useState<MatchOffer[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [deletingId, setDeletingId] = useState<string | null>(null)
  const [snackbar, setSnackbar] = useState<{ open: boolean; message: string; severity: 'success' | 'error' }>({ open: false, message: '', severity: 'success' })

  const [allRequests, setAllRequests] = useState<MatchRequest[]>([])
  const [requesterMap, setRequesterMap] = useState<Record<string, Account>>({})
  const [selectedOffer, setSelectedOffer] = useState<MatchOffer | null>(null)
  const [requestsDialogOpen, setRequestsDialogOpen] = useState(false)
  const [acceptingId, setAcceptingId] = useState<string | null>(null)

  useEffect(() => {
    findAccountMatchOffersUseCase.execute(accountId ?? '').then((result) => {
      if (result.success) setOffers(result.offers)
      else setError(result.error ?? 'Error al cargar las ofertas')
      setLoading(false)
    })

    findReceivedMatchRequestsUseCase.execute(accountId ?? '').then(async (result) => {
      if (!result.success) return
      setAllRequests(result.requests)

      const uniqueIds = [...new Set(result.requests.map((r) => r.requester_account_id))]
      const entries = await Promise.all(
        uniqueIds.map(async (id) => {
          try {
            const account = await fetchAccount(id)
            return [id, account] as [string, Account]
          } catch {
            return null
          }
        })
      )
      const map: Record<string, Account> = {}
      for (const entry of entries) {
        if (entry) map[entry[0]] = entry[1]
      }
      setRequesterMap(map)
    })
  }, [])

  const handleDelete = async (id: string) => {
    setDeletingId(id)
    const result = await deleteMatchOfferUseCase.execute(accountId ?? '', id)
    setDeletingId(null)
    if (result.success) {
      setOffers((prev) => prev.filter((o) => o.id !== id))
      setSnackbar({ open: true, message: 'Oferta eliminada correctamente', severity: 'success' })
    } else {
      setSnackbar({ open: true, message: result.error || 'Error al eliminar la oferta', severity: 'error' })
    }
  }

  const handleOpenRequests = (offer: MatchOffer) => {
    setSelectedOffer(offer)
    setRequestsDialogOpen(true)
  }

  const handleAccept = async (requestId: string) => {
    setAcceptingId(requestId)
    const result = await acceptMatchRequestUseCase.execute(accountId ?? '', requestId)
    setAcceptingId(null)
    if (result.success) {
      setAllRequests((prev) => prev.map((r) => r.id === requestId ? { ...r, status: 'ACCEPTED' } : r))
      setSnackbar({ open: true, message: 'Solicitud aceptada correctamente', severity: 'success' })
    } else {
      setSnackbar({ open: true, message: result.error ?? 'Error al aceptar la solicitud', severity: 'error' })
    }
  }

  const requestsForOffer = selectedOffer
    ? allRequests.filter((r) => r.match_offer_id === selectedOffer.id)
    : []

  const requestCountForOffer = (offerId: string) =>
    allRequests.filter((r) => r.match_offer_id === offerId).length

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('es-AR', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
  }

  const formatTime = (dateTimeString: string) => {
    const date = new Date(dateTimeString)
    const h = String(date.getHours()).padStart(2, '0')
    const m = String(date.getMinutes()).padStart(2, '0')
    return `${h}:${m}`
  }

  const getCategoryText = (admittedCategories: MatchOffer['admitted_categories']) => {
    switch (admittedCategories.type) {
      case 'SPECIFIC': return `Categorías: ${admittedCategories.categories?.map((c) => `L${c}`).join(', ')}`
      case 'GREATER_THAN': return `Nivel >= L${admittedCategories.min_level}`
      case 'LESS_THAN': return `Nivel <= L${admittedCategories.max_level}`
      case 'BETWEEN': return `Nivel L${admittedCategories.min_level} - L${admittedCategories.max_level}`
      default: return 'Cualquier nivel'
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'PENDING': return 'warning'
      case 'ACCEPTED': return 'success'
      case 'CANCELLED': case 'REJECTED': return 'error'
      default: return 'default'
    }
  }

  const getStatusText = (status: string) => {
    switch (status) {
      case 'PENDING': return 'Pendiente'
      case 'ACCEPTED': return 'Aceptada'
      case 'CANCELLED': return 'Cancelada'
      case 'REJECTED': return 'Rechazada'
      default: return status
    }
  }

  return (
    <>
      <Box>
        <Stack spacing={4}>
          <Paper elevation={0} sx={{ p: 4, background: 'linear-gradient(135deg, #6A1B9A 0%, #00C853 100%)', color: 'white', borderRadius: 4 }}>
            <Stack spacing={2} alignItems="center">
              <EventIcon sx={{ fontSize: 64 }} />
              <Typography variant="h3" component="h1" align="center" fontWeight={700}>Mis Publicaciones</Typography>
              <Typography variant="h6" align="center" sx={{ opacity: 0.95, maxWidth: 600 }}>Partidos que publicaste buscando rival</Typography>
              <Button
                variant="contained"
                startIcon={<AddIcon />}
                onClick={() => navigate('/my-offers/new')}
                sx={{ bgcolor: 'white', color: 'primary.main', '&:hover': { bgcolor: 'grey.100' } }}
              >
                Agregar publicación
              </Button>
            </Stack>
          </Paper>

          {loading && <Box display="flex" justifyContent="center" py={4}><CircularProgress /></Box>}
          {!loading && error && <Alert severity="error">{error}</Alert>}

          {!loading && !error && offers.length === 0 && (
            <Paper sx={{ p: 4, textAlign: 'center' }}>
              <EventIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
              <Typography variant="h6" color="text.secondary">No publicaste ninguna oferta todavía</Typography>
            </Paper>
          )}

          {!loading && !error && offers.length > 0 && (
            <Grid container spacing={3}>
              {offers.map((offer) => {
                const reqCount = requestCountForOffer(offer.id ?? '')
                return (
                  <Grid item xs={12} sm={6} md={4} key={offer.id}>
                    <Card elevation={2} sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                      <CardContent sx={{ flexGrow: 1 }}>
                        <Stack spacing={2}>
                          <Box display="flex" justifyContent="space-between" alignItems="start">
                            <Box>
                              <Typography variant="h6" fontWeight={700} gutterBottom>{offer.team_name}</Typography>
                              <Chip label={offer.sport} icon={<SportsIcon />} size="small" color="primary" variant="outlined" />
                            </Box>
                            <Chip label={getStatusText(offer.status)} size="small" color={getStatusColor(offer.status) as any} />
                          </Box>

                          <Divider />

                          <Box>
                            <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 1 }}>
                              <EventIcon fontSize="small" color="action" />
                              <Typography variant="body2" color="text.secondary" fontWeight={600}>{formatDate(offer.day)}</Typography>
                            </Stack>
                            <Stack direction="row" spacing={1} alignItems="center">
                              <AccessTimeIcon fontSize="small" color="action" />
                              <Typography variant="body2" color="text.secondary">
                                {formatTime(offer.time_slot.start_time)} - {formatTime(offer.time_slot.end_time)}
                              </Typography>
                            </Stack>
                          </Box>

                          <Stack direction="row" spacing={1} alignItems="center">
                            <LocationOnIcon fontSize="small" color="action" />
                            <Typography variant="body2" color="text.secondary">{offer.location.locality}, {offer.location.province}</Typography>
                          </Stack>

                          <Stack direction="row" spacing={1} alignItems="center">
                            <GroupsIcon fontSize="small" color="action" />
                            <Typography variant="body2" color="text.secondary">{getCategoryText(offer.admitted_categories)}</Typography>
                          </Stack>
                        </Stack>
                      </CardContent>

                      <CardActions sx={{ px: 2, pb: 2, gap: 1 }}>
                        <Button
                          variant="outlined"
                          fullWidth
                          startIcon={<InboxIcon />}
                          onClick={() => handleOpenRequests(offer)}
                        >
                          Solicitudes {reqCount > 0 && `(${reqCount})`}
                        </Button>
                        <Button
                          variant="outlined"
                          color="error"
                          fullWidth
                          startIcon={<DeleteIcon />}
                          disabled={deletingId === offer.id}
                          onClick={() => offer.id && handleDelete(offer.id)}
                        >
                          {deletingId === offer.id ? 'Eliminando...' : 'Eliminar'}
                        </Button>
                      </CardActions>
                    </Card>
                  </Grid>
                )
              })}
            </Grid>
          )}
        </Stack>
      </Box>

      {/* Requests Dialog */}
      <Dialog open={requestsDialogOpen} onClose={() => setRequestsDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          Solicitudes — {selectedOffer?.team_name}
        </DialogTitle>
        <DialogContent>
          {requestsForOffer.length === 0 ? (
            <Stack alignItems="center" py={3} spacing={1}>
              <InboxIcon sx={{ fontSize: 48, color: 'text.secondary' }} />
              <Typography color="text.secondary">No hay solicitudes para esta oferta</Typography>
            </Stack>
          ) : (
            <List disablePadding>
              {requestsForOffer.map((req, i) => {
                const requester = requesterMap[req.requester_account_id]
                const requesterName = requester
                  ? `${requester.FirstName} ${requester.LastName}`.trim() || requester.Nickname || requester.Email
                  : req.requester_account_id
                return (
                  <Box key={req.id}>
                    {i > 0 && <Divider />}
                    <ListItem sx={{ px: 0 }}>
                      <Stack spacing={1.5} width="100%">
                        <Stack direction="row" spacing={1.5} alignItems="center">
                          <Avatar src={requester?.Picture} alt={requesterName} sx={{ width: 40, height: 40, flexShrink: 0 }}>
                            {!requester?.Picture && <PersonIcon />}
                          </Avatar>
                          <Box>
                            <Typography variant="body2" fontWeight={600}>{requesterName}</Typography>
                            {requester?.Nickname && (
                              <Typography variant="caption" color="text.secondary">@{requester.Nickname}</Typography>
                            )}
                          </Box>
                        </Stack>

                        {selectedOffer && (
                          <Stack spacing={0.25} pl={7}>
                            <Stack direction="row" spacing={0.5} alignItems="center">
                              <EventIcon sx={{ fontSize: 13, color: 'text.secondary' }} />
                              <Typography variant="caption" color="text.secondary">{formatDate(selectedOffer.day)}</Typography>
                            </Stack>
                            <Stack direction="row" spacing={2}>
                              <Stack direction="row" spacing={0.5} alignItems="center">
                                <AccessTimeIcon sx={{ fontSize: 13, color: 'text.secondary' }} />
                                <Typography variant="caption" color="text.secondary">
                                  {formatTime(selectedOffer.time_slot.start_time)} - {formatTime(selectedOffer.time_slot.end_time)}
                                </Typography>
                              </Stack>
                              <Stack direction="row" spacing={0.5} alignItems="center">
                                <LocationOnIcon sx={{ fontSize: 13, color: 'text.secondary' }} />
                                <Typography variant="caption" color="text.secondary">
                                  {selectedOffer.location.locality}, {selectedOffer.location.province}
                                </Typography>
                              </Stack>
                            </Stack>
                          </Stack>
                        )}

                        <Box pl={7}>
                          {req.status === 'PENDING' ? (
                            <Stack direction="row" spacing={1}>
                              <Button
                                variant="contained"
                                color="success"
                                size="small"
                                startIcon={<CheckIcon />}
                                disabled={acceptingId === req.id}
                                onClick={() => handleAccept(req.id)}
                              >
                                {acceptingId === req.id ? 'Aceptando...' : 'Aceptar'}
                              </Button>
                              <Button variant="outlined" color="error" size="small" startIcon={<CloseIcon />} disabled>
                                Rechazar
                              </Button>
                            </Stack>
                          ) : (
                            <Chip label={getStatusText(req.status)} size="small" color={getStatusColor(req.status) as any} />
                          )}
                        </Box>
                      </Stack>
                    </ListItem>
                  </Box>
                )
              })}
            </List>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setRequestsDialogOpen(false)}>Cerrar</Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={snackbar.open}
        autoHideDuration={4000}
        onClose={() => setSnackbar((s) => ({ ...s, open: false }))}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert severity={snackbar.severity} onClose={() => setSnackbar((s) => ({ ...s, open: false }))}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </>
  )
}
