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
  Avatar,
  Card,
  CardContent,
} from '@mui/material'
import EmojiEventsIcon from '@mui/icons-material/EmojiEvents'
import EventIcon from '@mui/icons-material/Event'
import AccessTimeIcon from '@mui/icons-material/AccessTime'
import PersonIcon from '@mui/icons-material/Person'
import { useAuth } from '../../../auth/context/AuthContext'
import { fetchMatches } from '../../infrastructure/MatchApiAdapter'
import { fetchAccount, Account } from '../../../auth/infrastructure/adapters/AccountApiAdapter'
import { Match } from '../../domain/types'

function statusColor(status: string): 'success' | 'default' | 'error' {
  switch (status) {
    case 'ACCEPTED': return 'success'
    case 'CANCELLED': return 'error'
    default: return 'default'
  }
}

function statusText(status: string): string {
  switch (status) {
    case 'ACCEPTED': return 'Confirmado'
    case 'PLAYED': return 'Jugado'
    case 'CANCELLED': return 'Cancelado'
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

function accountDisplayName(acc: Account): string {
  return `${acc.FirstName} ${acc.LastName}`.trim() || acc.Nickname || acc.Email
}

export function MyMatchesPage() {
  const { accountId } = useAuth()

  const [matches, setMatches] = useState<Match[]>([])
  const [accountMap, setAccountMap] = useState<Record<string, Account>>({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!accountId) return

    fetchMatches(accountId)
      .then(async (result) => {
        setMatches(result)

        const allParticipantIds = [...new Set(result.flatMap((m) => m.participants))]

        const entries = await Promise.all(
          allParticipantIds.map(async (id) => {
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
        setAccountMap(map)
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }, [accountId])

  return (
    <Box>
      <Stack spacing={4}>
        <Paper
          elevation={0}
          sx={{ p: 4, background: 'linear-gradient(135deg, #1565C0 0%, #00C853 100%)', color: 'white', borderRadius: 4 }}
        >
          <Stack spacing={2} alignItems="center">
            <EmojiEventsIcon sx={{ fontSize: 64 }} />
            <Typography variant="h3" component="h1" align="center" fontWeight={700}>
              Mis Partidos
            </Typography>
            <Typography variant="h6" align="center" sx={{ opacity: 0.95, maxWidth: 600 }}>
              Todos tus partidos confirmados y jugados
            </Typography>
          </Stack>
        </Paper>

        {loading && <Box display="flex" justifyContent="center" py={4}><CircularProgress /></Box>}
        {!loading && error && <Alert severity="error">{error}</Alert>}

        {!loading && !error && matches.length === 0 && (
          <Paper sx={{ p: 4, textAlign: 'center' }}>
            <EmojiEventsIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" color="text.secondary">
              Todavía no tenés partidos registrados
            </Typography>
          </Paper>
        )}

        {!loading && !error && matches.length > 0 && (
          <Stack spacing={3}>
            {matches.map((match) => (
              <Card elevation={2} key={match.id}>
                  <CardContent>
                    <Stack spacing={2}>

                      {/* Header: título + estado */}
                      <Box display="flex" justifyContent="space-between" alignItems="flex-start">
                        <Typography variant="h6" fontWeight={700} sx={{ flex: 1, mr: 1 }}>
                          {match.title || match.sport}
                        </Typography>
                        <Chip
                          label={statusText(match.status)}
                          size="small"
                          color={statusColor(match.status)}
                        />
                      </Box>

                      <Divider />

                      {/* Fecha */}
                      <Stack direction="row" spacing={1} alignItems="center">
                        <EventIcon fontSize="small" color="action" />
                        <Typography variant="body2" color="text.secondary" fontWeight={600}>
                          {formatDate(match.day)}
                        </Typography>
                      </Stack>

                      {/* Horario */}
                      {match.time_slot?.start_time && (
                        <Stack direction="row" spacing={1} alignItems="center">
                          <AccessTimeIcon fontSize="small" color="action" />
                          <Typography variant="body2" color="text.secondary">
                            {formatTime(match.time_slot.start_time)} - {formatTime(match.time_slot.end_time)}
                          </Typography>
                        </Stack>
                      )}

                      <Divider />

                      {/* Participantes */}
                      <Stack spacing={1}>
                        <Typography variant="caption" color="text.secondary" fontWeight={600} textTransform="uppercase">
                          Participantes
                        </Typography>
                        {match.participants.map((pid) => {
                          const acc = accountMap[pid]
                          const name = acc ? accountDisplayName(acc) : pid
                          const isMe = pid === accountId
                          return (
                            <Stack key={pid} direction="row" spacing={1.5} alignItems="center">
                              <Avatar
                                src={acc?.Picture}
                                alt={name}
                                imgProps={{ referrerPolicy: 'no-referrer' }}
                                sx={{ width: 32, height: 32 }}
                              >
                                <PersonIcon fontSize="small" />
                              </Avatar>
                              <Typography variant="body2" fontWeight={isMe ? 600 : 400}>
                                {name}{isMe ? ' (vos)' : ''}
                              </Typography>
                            </Stack>
                          )
                        })}
                      </Stack>

                    </Stack>
                  </CardContent>
              </Card>
            ))}
          </Stack>
        )}
      </Stack>
    </Box>
  )
}
