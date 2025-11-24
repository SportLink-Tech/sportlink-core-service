import { useState } from 'react'
import {
  Box,
  Button,
  Card,
  CardContent,
  TextField,
  Typography,
  Chip,
  Stack,
  CircularProgress,
  Snackbar,
  Alert,
  Paper,
} from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import AddCircleIcon from '@mui/icons-material/AddCircle'
import CheckCircleIcon from '@mui/icons-material/CheckCircle'
import { SportSelect } from '../../../../shared/components/atoms/SportSelect'
import { CategorySelect } from '../../../../shared/components/atoms/CategorySelect'
import { Sport } from '../../../../shared/types/team'
import { useTeamContext } from '../../context/TeamContext'

const SPORTS: Sport[] = ['Football', 'Paddle']

/**
 * Feature Page: Create Team
 * Uses CreateTeamUseCase from context
 * Following Atomic Hexagonal Architecture
 */
export function CreateTeamPage() {
  const { createTeamUseCase } = useTeamContext()

  const [sport, setSport] = useState<Sport>('Paddle')
  const [name, setName] = useState('')
  const [category, setCategory] = useState<number>(0)
  const [playerIds, setPlayerIds] = useState<string[]>([])
  const [currentPlayerId, setCurrentPlayerId] = useState('')
  const [loading, setLoading] = useState(false)
  const [successTeamName, setSuccessTeamName] = useState<string | null>(null)
  const [showSuccessSnackbar, setShowSuccessSnackbar] = useState(false)
  const [showErrorSnackbar, setShowErrorSnackbar] = useState(false)

  const handleAddPlayer = () => {
    if (currentPlayerId.trim() && !playerIds.includes(currentPlayerId.trim())) {
      setPlayerIds([...playerIds, currentPlayerId.trim()])
      setCurrentPlayerId('')
    }
  }

  const handleRemovePlayer = (playerId: string) => {
    setPlayerIds(playerIds.filter(id => id !== playerId))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    // Execute use case
    const result = await createTeamUseCase.execute({
      sport,
      name,
      category,
      players: playerIds.length > 0 ? playerIds : undefined,
    })

    if (result.success && result.team) {
      setSuccessTeamName(result.team.Name)
      setShowSuccessSnackbar(true)
      // Reset form
      setName('')
      setCategory(0)
      setPlayerIds([])
      setCurrentPlayerId('')
    } else {
      setShowErrorSnackbar(true)
    }

    setLoading(false)
  }

  return (
    <>
      <Box>
        <Stack spacing={4}>
          {/* Hero Section */}
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
              <AddCircleIcon sx={{ fontSize: 64 }} />
              <Typography variant="h3" component="h1" align="center" fontWeight={700}>
                Crear Nuevo Equipo
              </Typography>
              <Typography variant="h6" align="center" sx={{ opacity: 0.95, maxWidth: 600 }}>
                Registra un nuevo equipo deportivo con su información, categoría y miembros
              </Typography>
            </Stack>
          </Paper>

          {/* Form Card */}
          <Card>
            <CardContent>
              <Typography variant="h5" component="h2" gutterBottom>
                Crear Nuevo Equipo
              </Typography>

              <Box component="form" onSubmit={handleSubmit}>
                <Stack spacing={3}>
                  <SportSelect
                    value={sport}
                    onChange={setSport}
                    required
                    sports={SPORTS}
                  />

                  <TextField
                    label="Nombre del Equipo"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    required
                    fullWidth
                    placeholder="Ej: Thunder Strikers"
                  />

                  <CategorySelect
                    value={category}
                    onChange={setCategory}
                  />

                  {/* Players Section */}
                  <Box>
                    <Typography variant="subtitle2" gutterBottom>
                      Jugadores (Opcional)
                    </Typography>
                    <Stack direction="row" spacing={1} sx={{ mb: 2 }}>
                      <TextField
                        size="small"
                        label="ID del Jugador"
                        value={currentPlayerId}
                        onChange={(e) => setCurrentPlayerId(e.target.value)}
                        placeholder="player-001"
                        fullWidth
                        onKeyPress={(e) => {
                          if (e.key === 'Enter') {
                            e.preventDefault()
                            handleAddPlayer()
                          }
                        }}
                      />
                      <Button
                        variant="outlined"
                        onClick={handleAddPlayer}
                        disabled={!currentPlayerId.trim()}
                        startIcon={<AddIcon />}
                      >
                        Agregar
                      </Button>
                    </Stack>

                    {playerIds.length > 0 && (
                      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                        {playerIds.map((id) => (
                          <Chip
                            key={id}
                            label={id}
                            onDelete={() => handleRemovePlayer(id)}
                            color="primary"
                            variant="outlined"
                          />
                        ))}
                      </Box>
                    )}
                  </Box>

                  <Button
                    type="submit"
                    variant="contained"
                    size="large"
                    disabled={loading || !name}
                    startIcon={loading ? <CircularProgress size={20} /> : <AddIcon />}
                  >
                    {loading ? 'Creando...' : 'Crear Equipo'}
                  </Button>
                </Stack>
              </Box>
            </CardContent>
          </Card>
        </Stack>
      </Box>

      {/* Success Snackbar */}
      <Snackbar
        open={showSuccessSnackbar}
        autoHideDuration={4000}
        onClose={() => setShowSuccessSnackbar(false)}
        anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
      >
        <Alert
          onClose={() => setShowSuccessSnackbar(false)}
          severity="success"
          variant="filled"
          icon={<CheckCircleIcon />}
          sx={{
            width: '100%',
            fontSize: '1.1rem',
            '& .MuiAlert-icon': {
              fontSize: '2rem',
            },
          }}
        >
          Equipo "{successTeamName}" creado exitosamente
        </Alert>
      </Snackbar>

      {/* Error Snackbar */}
      <Snackbar
        open={showErrorSnackbar}
        autoHideDuration={4000}
        onClose={() => setShowErrorSnackbar(false)}
        anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
      >
        <Alert
          onClose={() => setShowErrorSnackbar(false)}
          severity="error"
          variant="filled"
          sx={{
            width: '100%',
            fontSize: '1.1rem',
            '& .MuiAlert-icon': {
              fontSize: '2rem',
            },
          }}
        >
          Error al crear el equipo
        </Alert>
      </Snackbar>
    </>
  )
}

