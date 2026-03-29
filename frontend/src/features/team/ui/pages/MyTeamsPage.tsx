import { useEffect, useState } from 'react'
import {
  Box,
  Card,
  CardContent,
  Typography,
  Stack,
  CircularProgress,
  Paper,
  Chip,
} from '@mui/material'
import GroupsIcon from '@mui/icons-material/Groups'
import { useTeamContext } from '../../context/TeamContext'
import { Team } from '../../../../shared/types/team'
import { CURRENT_ACCOUNT_ID } from '../../../../shared/constants/session'

export function MyTeamsPage() {
  const { listAccountTeamsUseCase } = useTeamContext()
  const [teams, setTeams] = useState<Team[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    listAccountTeamsUseCase.execute(CURRENT_ACCOUNT_ID).then((result) => {
      if (result.success) {
        setTeams(result.teams)
      } else {
        setError(result.error ?? 'Error cargando equipos')
      }
      setLoading(false)
    })
  }, [])

  return (
    <Box>
      <Stack spacing={4}>
        <Paper
          elevation={0}
          sx={{
            p: 4,
            background: 'linear-gradient(135deg, #6A1B9A 0%, #00C853 100%)',
            color: 'white',
            borderRadius: 4,
          }}
        >
          <Stack spacing={2} alignItems="center">
            <GroupsIcon sx={{ fontSize: 64 }} />
            <Typography variant="h3" component="h1" align="center" fontWeight={700}>
              Mis Equipos
            </Typography>
            <Typography variant="h6" align="center" sx={{ opacity: 0.95, maxWidth: 600 }}>
              Equipos asociados a tu cuenta
            </Typography>
          </Stack>
        </Paper>

        {loading && (
          <Box display="flex" justifyContent="center" py={4}>
            <CircularProgress />
          </Box>
        )}

        {!loading && error && (
          <Typography color="error" align="center">{error}</Typography>
        )}

        {!loading && !error && teams.length === 0 && (
          <Typography align="center" color="text.secondary">
            No tenés equipos registrados todavía.
          </Typography>
        )}

        {!loading && !error && teams.length > 0 && (
          <Stack spacing={2}>
            {teams.map((team) => (
              <Card key={team.Name + team.Sport}>
                <CardContent>
                  <Stack direction="row" alignItems="center" justifyContent="space-between">
                    <Box>
                      <Typography variant="h6" fontWeight={600}>{team.Name}</Typography>
                      <Typography variant="body2" color="text.secondary">{team.Sport}</Typography>
                    </Box>
                    <Chip label={`Categoría ${team.Category}`} color="primary" variant="outlined" />
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
