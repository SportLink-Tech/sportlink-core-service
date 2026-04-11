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
} from '@mui/material'
import EmojiEventsIcon from '@mui/icons-material/EmojiEvents'
import EventIcon from '@mui/icons-material/Event'
import SportsIcon from '@mui/icons-material/Sports'
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

        const opponentIds = result.map((m) =>
          m.local_account_id === accountId ? m.visitor_account_id : m.local_account_id
        )
        const uniqueIds = [...new Set(opponentIds)]

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
        setAccountMap(map)
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }, [accountId])

  const opponentOf = (m: Match): Account | undefined => {
    const opponentId = m.local_account_id === accountId ? m.visitor_account_id : m.local_account_id
    return accountMap[opponentId]
  }

  const opponentName = (m: Match): string => {
    const acc = opponentOf(m)
    const opponentId = m.local_account_id === accountId ? m.visitor_account_id : m.local_account_id
    if (!acc) return opponentId
    return `${acc.FirstName} ${acc.LastName}`.trim() || acc.Nickname || acc.Email
  }

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
          <Paper elevation={0} sx={{ borderRadius: 2, border: '1px solid', borderColor: 'divider' }}>
            <List disablePadding>
              {matches.map((match, i) => {
                const opponent = opponentOf(match)
                const name = opponentName(match)
                const isLocal = match.local_account_id === accountId
                return (
                  <Box key={match.id}>
                    {i > 0 && <Divider />}
                    <ListItem sx={{ px: 3, py: 2 }}>
                      <Stack spacing={1} width="100%">
                        <Stack direction="row" spacing={2} alignItems="center" justifyContent="space-between">
                          <Stack direction="row" spacing={2} alignItems="center">
                            <Avatar src={opponent?.Picture} alt={name} sx={{ width: 40, height: 40 }}>
                              {!opponent?.Picture && <PersonIcon />}
                            </Avatar>
                            <Box>
                              <Typography variant="body1" fontWeight={600}>{name}</Typography>
                              {opponent?.Nickname && (
                                <Typography variant="caption" color="text.secondary">@{opponent.Nickname}</Typography>
                              )}
                            </Box>
                          </Stack>
                          <Chip label={statusText(match.status)} size="small" color={statusColor(match.status)} />
                        </Stack>

                        <Stack direction="row" spacing={3} pl={7}>
                          <Stack direction="row" spacing={0.5} alignItems="center">
                            <EventIcon sx={{ fontSize: 14, color: 'text.secondary' }} />
                            <Typography variant="caption" color="text.secondary">
                              {formatDate(match.day)}
                            </Typography>
                          </Stack>
                          <Stack direction="row" spacing={0.5} alignItems="center">
                            <SportsIcon sx={{ fontSize: 14, color: 'text.secondary' }} />
                            <Typography variant="caption" color="text.secondary">{match.sport}</Typography>
                          </Stack>
                          <Typography variant="caption" color="text.secondary">
                            {isLocal ? 'Local' : 'Visitante'}
                          </Typography>
                        </Stack>
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
  )
}
