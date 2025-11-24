import { useState } from 'react'
import {
  Box,
  Button,
  Card,
  CardContent,
  TextField,
  MenuItem,
  Typography,
  Alert,
  Chip,
  Stack,
  CircularProgress,
  Snackbar,
} from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import CheckCircleIcon from '@mui/icons-material/CheckCircle'
import { apiService } from '../services/api'
import { Sport, Team } from '../types/team'

const SPORTS: Sport[] = ['Football', 'Paddle']
const CATEGORIES = [
  { value: 0, label: 'Unranked' },
  { value: 1, label: 'L1 - Principiante' },
  { value: 2, label: 'L2' },
  { value: 3, label: 'L3' },
  { value: 4, label: 'L4' },
  { value: 5, label: 'L5' },
  { value: 6, label: 'L6' },
  { value: 7, label: 'L7 - Avanzado' },
]

export function CreateTeamForm() {
  const [sport, setSport] = useState<Sport>('Paddle')
  const [name, setName] = useState('')
  const [category, setCategory] = useState<number>(0)
  const [playerIds, setPlayerIds] = useState<string[]>([])
  const [currentPlayerId, setCurrentPlayerId] = useState('')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState<Team | null>(null)
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
    setSuccess(null)
    setLoading(true)

    try {
      const response = await apiService.createTeam({
        sport,
        name,
        category,
        players: playerIds.length > 0 ? playerIds : undefined,
      })
      
      if (response.status === 201) {
        setSuccess(response.data)
        setShowSuccessSnackbar(true)
        // Reset form
        setName('')
        setCategory(0)
        setPlayerIds([])
        setCurrentPlayerId('')
      } else {
        setShowErrorSnackbar(true)
      }
    } catch (err) {
      setShowErrorSnackbar(true)
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <Card>
        <CardContent>
          <Typography variant="h5" component="h2" gutterBottom>
            Crear Nuevo Equipo
          </Typography>

        <Box component="form" onSubmit={handleSubmit}>
          <Stack spacing={3}>
            {/* Sport Selection */}
            <TextField
              select
              label="Deporte"
              value={sport}
              onChange={(e) => setSport(e.target.value as Sport)}
              required
              fullWidth
            >
              {SPORTS.map((s) => (
                <MenuItem key={s} value={s}>
                  {s}
                </MenuItem>
              ))}
            </TextField>

            {/* Team Name */}
            <TextField
              label="Nombre del Equipo"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              fullWidth
              placeholder="Ej: Thunder Strikers"
            />

            {/* Category */}
            <TextField
              select
              label="Categoría"
              value={category}
              onChange={(e) => setCategory(Number(e.target.value))}
              fullWidth
            >
              {CATEGORIES.map((cat) => (
                <MenuItem key={cat.value} value={cat.value}>
                  {cat.label}
                </MenuItem>
              ))}
            </TextField>

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

            {/* Submit Button */}
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
        Equipo "{success?.Name}" creado exitosamente
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

