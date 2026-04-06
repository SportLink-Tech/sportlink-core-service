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
  IconButton,
  TextField,
  Tooltip,
} from '@mui/material'
import GroupsIcon from '@mui/icons-material/Groups'
import EditIcon from '@mui/icons-material/Edit'
import CheckIcon from '@mui/icons-material/Check'
import CloseIcon from '@mui/icons-material/Close'
import { useTeamContext } from '../../context/TeamContext'
import { Team } from '../../../../shared/types/team'
import { useAuth } from '../../../auth/context/AuthContext'

export function MyTeamsPage() {
  const { listAccountTeamsUseCase, updateTeamUseCase } = useTeamContext()
  const { accountId } = useAuth()
  const [teams, setTeams] = useState<Team[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [editingKey, setEditingKey] = useState<string | null>(null)
  const [editName, setEditName] = useState('')
  const [saveError, setSaveError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    listAccountTeamsUseCase.execute(accountId ?? '').then((result) => {
      if (result.success) {
        setTeams(result.teams)
      } else {
        setError(result.error ?? 'Error cargando equipos')
      }
      setLoading(false)
    })
  }, [])

  const teamKey = (team: Team) => `${team.Sport}#${team.Name}`

  const startEditing = (team: Team) => {
    setEditingKey(teamKey(team))
    setEditName(team.Name)
    setSaveError(null)
  }

  const cancelEditing = () => {
    setEditingKey(null)
    setEditName('')
    setSaveError(null)
  }

  const saveTeam = async (team: Team) => {
    if (!editName.trim()) {
      setSaveError('El nombre no puede estar vacío')
      return
    }

    setSaving(true)
    setSaveError(null)

    const result = await updateTeamUseCase.execute(team.Sport, team.Name, { name: editName.trim() })

    setSaving(false)

    if (result.success && result.team) {
      setTeams((prev) =>
        prev.map((t) => (teamKey(t) === teamKey(team) ? result.team! : t))
      )
      setEditingKey(null)
    } else {
      setSaveError(result.error ?? 'Error actualizando el equipo')
    }
  }

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
            {teams.map((team) => {
              const key = teamKey(team)
              const isEditing = editingKey === key

              return (
                <Card key={key}>
                  <CardContent>
                    <Stack direction="row" alignItems="center" justifyContent="space-between" spacing={2}>
                      <Box flex={1}>
                        {isEditing ? (
                          <Stack spacing={1}>
                            <TextField
                              size="small"
                              value={editName}
                              onChange={(e) => setEditName(e.target.value)}
                              onKeyDown={(e) => {
                                if (e.key === 'Enter') saveTeam(team)
                                if (e.key === 'Escape') cancelEditing()
                              }}
                              error={!!saveError}
                              helperText={saveError ?? undefined}
                              autoFocus
                              fullWidth
                            />
                          </Stack>
                        ) : (
                          <>
                            <Typography variant="h6" fontWeight={600}>{team.Name}</Typography>
                            <Typography variant="body2" color="text.secondary">{team.Sport}</Typography>
                          </>
                        )}
                      </Box>

                      <Stack direction="row" alignItems="center" spacing={1}>
                        {!isEditing && (
                          <Chip label={`Categoría ${team.Category}`} color="primary" variant="outlined" />
                        )}

                        {isEditing ? (
                          <>
                            <Tooltip title="Guardar">
                              <span>
                                <IconButton
                                  color="primary"
                                  onClick={() => saveTeam(team)}
                                  disabled={saving}
                                  size="small"
                                >
                                  <CheckIcon />
                                </IconButton>
                              </span>
                            </Tooltip>
                            <Tooltip title="Cancelar">
                              <IconButton onClick={cancelEditing} size="small">
                                <CloseIcon />
                              </IconButton>
                            </Tooltip>
                          </>
                        ) : (
                          <Tooltip title="Editar nombre">
                            <IconButton onClick={() => startEditing(team)} size="small">
                              <EditIcon />
                            </IconButton>
                          </Tooltip>
                        )}
                      </Stack>
                    </Stack>
                  </CardContent>
                </Card>
              )
            })}
          </Stack>
        )}
      </Stack>
    </Box>
  )
}
