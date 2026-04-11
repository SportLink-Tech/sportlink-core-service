import { useEffect, useState } from 'react'
import {
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Divider,
  IconButton,
  Paper,
  Stack,
  TextField,
  Tooltip,
  Typography,
  Snackbar,
  Alert,
} from '@mui/material'
import GroupsIcon from '@mui/icons-material/Groups'
import AddIcon from '@mui/icons-material/Add'
import EditIcon from '@mui/icons-material/Edit'
import CheckIcon from '@mui/icons-material/Check'
import CloseIcon from '@mui/icons-material/Close'
import CheckCircleIcon from '@mui/icons-material/CheckCircle'
import { useTeamContext } from '../../context/TeamContext'
import { useAuth } from '../../../auth/context/AuthContext'
import { Team, Sport } from '../../../../shared/types/team'
import { SportSelect } from '../../../../shared/components/atoms/SportSelect'
import { CategorySelect } from '../../../../shared/components/atoms/CategorySelect'

const SPORTS: Sport[] = ['Football', 'Paddle']

export function TeamsPage() {
  const { listAccountTeamsUseCase, createTeamUseCase, updateTeamUseCase } = useTeamContext()
  const { accountId } = useAuth()

  // List state
  const [teams, setTeams] = useState<Team[]>([])
  const [loading, setLoading] = useState(true)
  const [listError, setListError] = useState<string | null>(null)

  // Edit state
  const [editingKey, setEditingKey] = useState<string | null>(null)
  const [editName, setEditName] = useState('')
  const [editError, setEditError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  // Create state
  const [sport, setSport] = useState<Sport>('Paddle')
  const [name, setName] = useState('')
  const [category, setCategory] = useState<number>(0)
  const [creating, setCreating] = useState(false)
  const [successTeamName, setSuccessTeamName] = useState<string | null>(null)
  const [showSuccess, setShowSuccess] = useState(false)
  const [showCreateError, setShowCreateError] = useState(false)

  useEffect(() => {
    listAccountTeamsUseCase.execute(accountId ?? '').then((result) => {
      if (result.success) {
        setTeams(result.teams)
      } else {
        setListError(result.error ?? 'Error cargando equipos')
      }
      setLoading(false)
    })
  }, [])

  const teamKey = (team: Team) => `${team.Sport}#${team.Name}`

  // ── Edit ──────────────────────────────────────────────────────────────────

  const startEditing = (team: Team) => {
    setEditingKey(teamKey(team))
    setEditName(team.Name)
    setEditError(null)
  }

  const cancelEditing = () => {
    setEditingKey(null)
    setEditName('')
    setEditError(null)
  }

  const saveTeam = async (team: Team) => {
    if (!editName.trim()) {
      setEditError('El nombre no puede estar vacío')
      return
    }
    setSaving(true)
    setEditError(null)
    const result = await updateTeamUseCase.execute(team.Sport, team.Name, { name: editName.trim() })
    setSaving(false)
    if (result.success && result.team) {
      setTeams((prev) => prev.map((t) => (teamKey(t) === teamKey(team) ? result.team! : t)))
      setEditingKey(null)
    } else {
      setEditError(result.error ?? 'Error actualizando el equipo')
    }
  }

  // ── Create ────────────────────────────────────────────────────────────────

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setCreating(true)
    const result = await createTeamUseCase.execute(accountId ?? '', {
      sport,
      name,
      category,
      owner_account_id: accountId ?? '',
    })
    setCreating(false)
    if (result.success && result.team) {
      setTeams((prev) => [...prev, result.team!])
      setSuccessTeamName(result.team.Name)
      setShowSuccess(true)
      setName('')
      setCategory(0)
    } else {
      setShowCreateError(true)
    }
  }

  return (
    <>
      <Stack spacing={4}>
        {/* Header */}
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
              Gestioná tus equipos y creá nuevos desde acá
            </Typography>
          </Stack>
        </Paper>

        {/* Team list */}
        {loading && (
          <Box display="flex" justifyContent="center" py={2}>
            <CircularProgress />
          </Box>
        )}

        {!loading && listError && (
          <Alert severity="error">{listError}</Alert>
        )}

        {!loading && !listError && teams.length > 0 && (
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
                          <TextField
                            size="small"
                            value={editName}
                            onChange={(e) => setEditName(e.target.value)}
                            onKeyDown={(e) => {
                              if (e.key === 'Enter') saveTeam(team)
                              if (e.key === 'Escape') cancelEditing()
                            }}
                            error={!!editError}
                            helperText={editError ?? undefined}
                            autoFocus
                            fullWidth
                          />
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
                                <IconButton color="primary" onClick={() => saveTeam(team)} disabled={saving} size="small">
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

        {/* Create form */}
        <Card>
          <CardContent>
            {teams.length > 0 && <Divider sx={{ mb: 3 }} />}
            <Typography variant="h6" fontWeight={600} gutterBottom>
              Crear nuevo equipo
            </Typography>
            <Box component="form" onSubmit={handleCreate}>
              <Stack spacing={3}>
                <SportSelect value={sport} onChange={setSport} required sports={SPORTS} />
                <TextField
                  label="Nombre del equipo"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
                  fullWidth
                  placeholder="Ej: Thunder Strikers"
                />
                <CategorySelect value={category} onChange={setCategory} />
                <Button
                  type="submit"
                  variant="contained"
                  size="large"
                  disabled={creating || !name}
                  startIcon={creating ? <CircularProgress size={20} /> : <AddIcon />}
                >
                  {creating ? 'Creando...' : 'Crear equipo'}
                </Button>
              </Stack>
            </Box>
          </CardContent>
        </Card>
      </Stack>

      <Snackbar
        open={showSuccess}
        autoHideDuration={4000}
        onClose={() => setShowSuccess(false)}
        anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
      >
        <Alert onClose={() => setShowSuccess(false)} severity="success" variant="filled" icon={<CheckCircleIcon />}>
          Equipo "{successTeamName}" creado exitosamente
        </Alert>
      </Snackbar>

      <Snackbar
        open={showCreateError}
        autoHideDuration={4000}
        onClose={() => setShowCreateError(false)}
        anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
      >
        <Alert onClose={() => setShowCreateError(false)} severity="error" variant="filled">
          Error al crear el equipo
        </Alert>
      </Snackbar>
    </>
  )
}
