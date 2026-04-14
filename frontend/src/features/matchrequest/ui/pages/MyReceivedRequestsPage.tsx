import { useState, useEffect } from 'react'
import {
  Box,
  Typography,
  Stack,
  Paper,
  Chip,
  CircularProgress,
  Alert,
  Divider,
  List,
  ListItem,
  Avatar,
  Button,
  Snackbar,
} from '@mui/material'
import InboxIcon from '@mui/icons-material/Inbox'
import EventIcon from '@mui/icons-material/Event'
import AccessTimeIcon from '@mui/icons-material/AccessTime'
import PersonIcon from '@mui/icons-material/Person'
import LocationOnIcon from '@mui/icons-material/LocationOn'
import CheckIcon from '@mui/icons-material/Check'
import CloseIcon from '@mui/icons-material/Close'
import { useMatchRequestContext } from '../../context/MatchRequestContext'
import { useMatchOfferContext } from '../../../matchoffer/context/MatchOfferContext'
import { MatchRequest } from '../../domain/ports/MatchRequestRepository'
import { MatchOffer } from '../../../../shared/types/matchOffer'
import { useAuth } from '../../../auth/context/AuthContext'
import { fetchAccount, Account } from '../../../auth/infrastructure/adapters/AccountApiAdapter'

function statusColor(status: string): 'warning' | 'success' | 'error' | 'default' {
  switch (status) {
    case 'PENDING': return 'warning'
    case 'ACCEPTED': return 'success'
    case 'REJECTED': return 'error'
    default: return 'default'
  }
}

function statusText(status: string): string {
  switch (status) {
    case 'PENDING': return 'Pendiente'
    case 'ACCEPTED': return 'Aceptada'
    case 'REJECTED': return 'Rechazada'
    default: return status
  }
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString('es-AR', {
    weekday: 'long', day: '2-digit', month: 'long', year: 'numeric',
  })
}

function formatTime(dateTimeString: string): string {
  const d = new Date(dateTimeString)
  return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

export function MyReceivedRequestsPage() {
  const { findReceivedMatchRequestsUseCase, acceptMatchRequestUseCase } = useMatchRequestContext()
  const { retrieveMatchOfferUseCase } = useMatchOfferContext()
  const { accountId } = useAuth()

  const [requests, setRequests] = useState<MatchRequest[]>([])
  const [offerMap, setOfferMap] = useState<Record<string, MatchOffer>>({})
  const [requesterMap, setRequesterMap] = useState<Record<string, Account>>({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [acceptingId, setAcceptingId] = useState<string | null>(null)
  const [snackbar, setSnackbar] = useState<{ open: boolean; message: string; severity: 'success' | 'error' }>({ open: false, message: '', severity: 'success' })

  useEffect(() => {
    findReceivedMatchRequestsUseCase.execute(accountId ?? '', ['PENDING', 'ACCEPTED', 'REJECTED']).then(async (result) => {
      if (!result.success) {
        setError(result.error ?? 'Error al cargar las solicitudes')
        setLoading(false)
        return
      }
      setRequests(result.requests)

      const uniqueOfferIds = [...new Set(result.requests.map((r) => r.match_offer_id))]
      const uniqueRequesterIds = [...new Set(result.requests.map((r) => r.requester_account_id))]

      const [offerEntries, requesterEntries] = await Promise.all([
        Promise.all(
          uniqueOfferIds.map(async (id) => {
            const offer = await retrieveMatchOfferUseCase.execute(id)
            return offer ? ([id, offer] as [string, MatchOffer]) : null
          }),
        ),
        Promise.all(
          uniqueRequesterIds.map(async (id) => {
            try {
              const account = await fetchAccount(id)
              return [id, account] as [string, Account]
            } catch {
              return null
            }
          }),
        ),
      ])

      const offerMap: Record<string, MatchOffer> = {}
      for (const entry of offerEntries) {
        if (entry) offerMap[entry[0]] = entry[1]
      }

      const requesterMap: Record<string, Account> = {}
      for (const entry of requesterEntries) {
        if (entry) requesterMap[entry[0]] = entry[1]
      }

      setOfferMap(offerMap)
      setRequesterMap(requesterMap)
      setLoading(false)
    })
  }, [])

  const handleAccept = async (requestId: string) => {
    setAcceptingId(requestId)
    const result = await acceptMatchRequestUseCase.execute(accountId ?? '', requestId)
    setAcceptingId(null)
    if (result.success) {
      setRequests((prev) => prev.map((r) => r.id === requestId ? { ...r, status: 'ACCEPTED' } : r))
      setSnackbar({ open: true, message: 'Solicitud aceptada correctamente', severity: 'success' })
    } else {
      setSnackbar({ open: true, message: result.error ?? 'Error al aceptar la solicitud', severity: 'error' })
    }
  }

  return (
    <>
    <Box>
      <Stack spacing={4}>
        <Paper
          elevation={0}
          sx={{ p: 4, background: 'linear-gradient(135deg, #6A1B9A 0%, #E65100 100%)', color: 'white', borderRadius: 4 }}
        >
          <Stack spacing={2} alignItems="center">
            <InboxIcon sx={{ fontSize: 64 }} />
            <Typography variant="h3" component="h1" align="center" fontWeight={700}>
              Solicitudes Recibidas
            </Typography>
            <Typography variant="h6" align="center" sx={{ opacity: 0.95, maxWidth: 600 }}>
              Equipos que quieren jugar contra vos
            </Typography>
          </Stack>
        </Paper>

        {loading && <Box display="flex" justifyContent="center" py={4}><CircularProgress /></Box>}
        {!loading && error && <Alert severity="error">{error}</Alert>}

        {!loading && !error && requests.length === 0 && (
          <Paper sx={{ p: 4, textAlign: 'center' }}>
            <InboxIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" color="text.secondary">
              No recibiste ninguna solicitud todavía
            </Typography>
          </Paper>
        )}

        {!loading && !error && requests.length > 0 && (
          <Paper elevation={0} sx={{ borderRadius: 2, border: '1px solid', borderColor: 'divider' }}>
            <List disablePadding>
              {requests.map((req, i) => {
                const offer = offerMap[req.match_offer_id]
                const requester = requesterMap[req.requester_account_id]
                const requesterName = requester
                  ? `${requester.FirstName} ${requester.LastName}`.trim() || requester.Nickname || requester.Email
                  : req.requester_account_id
                return (
                  <Box key={req.id}>
                    {i > 0 && <Divider />}
                    <ListItem sx={{ px: 3, py: 2 }}>
                      <Stack spacing={1.5} width="100%">
                        <Stack direction="row" spacing={2} alignItems="center">
                          <Avatar src={requester?.Picture} alt={requesterName} sx={{ width: 40, height: 40, flexShrink: 0 }}>
                            <PersonIcon />
                          </Avatar>
                          <Box>
                            <Typography variant="body1" fontWeight={600}>{requesterName}</Typography>
                            {requester?.Nickname && (
                              <Typography variant="caption" color="text.secondary">@{requester.Nickname}</Typography>
                            )}
                          </Box>
                        </Stack>

                        {offer && (
                          <Stack spacing={0.5} pl={7}>
                            <Stack direction="row" spacing={2}>
                              <Stack direction="row" spacing={0.5} alignItems="center">
                                <EventIcon sx={{ fontSize: 14, color: 'text.secondary' }} />
                                <Typography variant="caption" color="text.secondary">{formatDate(offer.day)}</Typography>
                              </Stack>
                              <Stack direction="row" spacing={0.5} alignItems="center">
                                <AccessTimeIcon sx={{ fontSize: 14, color: 'text.secondary' }} />
                                <Typography variant="caption" color="text.secondary">
                                  {formatTime(offer.time_slot.start_time)} - {formatTime(offer.time_slot.end_time)}
                                </Typography>
                              </Stack>
                            </Stack>
                            <Stack direction="row" spacing={0.5} alignItems="center">
                              <LocationOnIcon sx={{ fontSize: 14, color: 'text.secondary' }} />
                              <Typography variant="caption" color="text.secondary">
                                {offer.location.locality}, {offer.location.province}
                              </Typography>
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
                            <Chip label={statusText(req.status)} size="small" color={statusColor(req.status)} />
                          )}
                        </Box>
                      </Stack>
                    </ListItem>
                  </Box>
                )
              })}
            </List>
          </Paper>
        )}
      </Stack>
    </Box>

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
