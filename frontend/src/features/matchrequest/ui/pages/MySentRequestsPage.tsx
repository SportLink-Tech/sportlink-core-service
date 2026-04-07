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
  ListItemText,
} from '@mui/material'
import SendIcon from '@mui/icons-material/Send'
import InboxIcon from '@mui/icons-material/Inbox'
import EventIcon from '@mui/icons-material/Event'
import AccessTimeIcon from '@mui/icons-material/AccessTime'
import { useMatchRequestContext } from '../../context/MatchRequestContext'
import { useMatchOfferContext } from '../../../matchoffer/context/MatchOfferContext'
import { MatchRequest } from '../../domain/ports/MatchRequestRepository'
import { MatchOffer } from '../../../../shared/types/matchOffer'
import { useAuth } from '../../../auth/context/AuthContext'

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

export function MySentRequestsPage() {
  const { findSentMatchRequestsListUseCase } = useMatchRequestContext()
  const { retrieveMatchOfferUseCase } = useMatchOfferContext()
  const { accountId } = useAuth()

  const [requests, setRequests] = useState<MatchRequest[]>([])
  const [offerMap, setOfferMap] = useState<Record<string, MatchOffer>>({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    findSentMatchRequestsListUseCase.execute(accountId ?? '').then(async (result) => {
      if (!result.success) {
        setError(result.error ?? 'Error al cargar las solicitudes')
        setLoading(false)
        return
      }
      setRequests(result.requests)

      const uniqueOfferIds = [...new Set(result.requests.map((r) => r.match_offer_id))]
      const entries = await Promise.all(
        uniqueOfferIds.map(async (id) => {
          const offer = await retrieveMatchOfferUseCase.execute(id)
          return offer ? ([id, offer] as [string, MatchOffer]) : null
        }),
      )
      const map: Record<string, MatchOffer> = {}
      for (const entry of entries) {
        if (entry) map[entry[0]] = entry[1]
      }
      setOfferMap(map)
      setLoading(false)
    })
  }, [])

  return (
    <Box>
      <Stack spacing={4}>
        <Paper
          elevation={0}
          sx={{ p: 4, background: 'linear-gradient(135deg, #1565C0 0%, #00C853 100%)', color: 'white', borderRadius: 4 }}
        >
          <Stack spacing={2} alignItems="center">
            <SendIcon sx={{ fontSize: 64 }} />
            <Typography variant="h3" component="h1" align="center" fontWeight={700}>
              Solicitudes Enviadas
            </Typography>
            <Typography variant="h6" align="center" sx={{ opacity: 0.95, maxWidth: 600 }}>
              Partidos a los que solicitaste unirte
            </Typography>
          </Stack>
        </Paper>

        {loading && <Box display="flex" justifyContent="center" py={4}><CircularProgress /></Box>}
        {!loading && error && <Alert severity="error">{error}</Alert>}

        {!loading && !error && requests.length === 0 && (
          <Paper sx={{ p: 4, textAlign: 'center' }}>
            <InboxIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" color="text.secondary">
              No enviaste ninguna solicitud todavía
            </Typography>
          </Paper>
        )}

        {!loading && !error && requests.length > 0 && (
          <Paper elevation={0} sx={{ borderRadius: 2, border: '1px solid', borderColor: 'divider' }}>
            <List disablePadding>
              {requests.map((req, i) => {
                const offer = offerMap[req.match_offer_id]
                return (
                  <Box key={req.id}>
                    {i > 0 && <Divider />}
                    <ListItem sx={{ px: 3, py: 2 }}>
                      <ListItemText
                        primary={
                          <Typography variant="body1" fontWeight={600}>
                            {offer ? offer.team_name : req.match_offer_id}
                          </Typography>
                        }
                        secondary={
                          offer ? (
                            <Stack direction="row" spacing={2} mt={0.5}>
                              <Stack direction="row" spacing={0.5} alignItems="center">
                                <EventIcon sx={{ fontSize: 14, color: 'text.secondary' }} />
                                <Typography variant="caption" color="text.secondary">
                                  {formatDate(offer.day)}
                                </Typography>
                              </Stack>
                              <Stack direction="row" spacing={0.5} alignItems="center">
                                <AccessTimeIcon sx={{ fontSize: 14, color: 'text.secondary' }} />
                                <Typography variant="caption" color="text.secondary">
                                  {formatTime(offer.time_slot.start_time)} - {formatTime(offer.time_slot.end_time)}
                                </Typography>
                              </Stack>
                            </Stack>
                          ) : null
                        }
                      />
                      <Chip label={statusText(req.status)} size="small" color={statusColor(req.status)} />
                    </ListItem>
                  </Box>
                )
              })}
            </List>
          </Paper>
        )}
      </Stack>
    </Box>
  )
}
