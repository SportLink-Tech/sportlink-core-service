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
  Stack,
  CircularProgress,
  Divider,
  Chip,
} from '@mui/material'
import SearchIcon from '@mui/icons-material/Search'
import { apiService } from '../services/api'
import { Sport, Team } from '../types/team'

const SPORTS: Sport[] = ['Football', 'Paddle', 'Tennis']

export function SearchTeamForm() {
  const [sport, setSport] = useState<Sport>('Paddle')
  const [teamName, setTeamName] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [teams, setTeams] = useState<Team[]>([])

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setTeams([])
    setLoading(true)

    try {
      const response = await apiService.findTeams(sport, teamName || undefined)
      console.log('Teams search response:', response)
      
      // Normalizar Members a array vacío si es null para cada team
      const teamsData = response.data.map(team => ({
        ...team,
        Members: team.Members || []
      }))
      
      setTeams(teamsData)
      
      if (teamsData.length === 0) {
        setError('No se encontraron equipos con los criterios especificados')
      }
    } catch (err) {
      console.error('Error searching teams:', err)
      setError(err instanceof Error ? err.message : 'Error al buscar equipos')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h5" component="h2" gutterBottom>
          Buscar Equipo
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        <Box component="form" onSubmit={handleSearch}>
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
              label="Nombre del Equipo (opcional)"
              value={teamName}
              onChange={(e) => setTeamName(e.target.value)}
              fullWidth
              placeholder="Ej: Thunder Strikers"
            />

            {/* Search Button */}
            <Button
              type="submit"
              variant="contained"
              size="large"
              disabled={loading}
              startIcon={loading ? <CircularProgress size={20} /> : <SearchIcon />}
            >
              {loading ? 'Buscando...' : 'Buscar Equipos'}
            </Button>
          </Stack>
        </Box>

        {/* Teams Results */}
        {teams.length > 0 && (
          <Box sx={{ mt: 4 }}>
            <Divider sx={{ mb: 3 }} />
            <Typography variant="h6" gutterBottom color="primary">
              ✓ {teams.length} {teams.length === 1 ? 'Equipo Encontrado' : 'Equipos Encontrados'}
            </Typography>
            
            <Stack spacing={3}>
              {teams.map((team, index) => (
                <Card key={`${team.Sport}-${team.Name}-${index}`} variant="outlined" sx={{ p: 2 }}>
            <Stack spacing={2}>
              <Box>
                <Typography variant="subtitle2" color="text.secondary">
                  Nombre
                </Typography>
                <Typography variant="h5">{team.Name}</Typography>
              </Box>

              <Box>
                <Typography variant="subtitle2" color="text.secondary">
                  Deporte
                </Typography>
                <Chip label={team.Sport} color="primary" />
              </Box>

              <Box>
                <Typography variant="subtitle2" color="text.secondary">
                  Categoría
                </Typography>
                <Chip 
                  label={team.Category === 0 ? 'Unranked' : `L${team.Category}`} 
                  color="secondary" 
                />
              </Box>

              <Box>
                <Typography variant="subtitle2" color="text.secondary">
                  Estadísticas
                </Typography>
                <Stack direction="row" spacing={2}>
                  <Chip label={`Victorias: ${team.Stats.Wins}`} variant="outlined" color="success" />
                  <Chip label={`Derrotas: ${team.Stats.Losses}`} variant="outlined" color="error" />
                  <Chip label={`Empates: ${team.Stats.Draws}`} variant="outlined" />
                </Stack>
              </Box>

              <Box>
                <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                  Miembros ({team.Members?.length || 0})
                </Typography>
                {team.Members && team.Members.length > 0 ? (
                  <Stack spacing={1}>
                    {team.Members.map((member) => (
                      <Card key={member.ID} variant="outlined">
                        <CardContent sx={{ py: 1 }}>
                          <Stack direction="row" spacing={2} alignItems="center">
                            <Typography variant="body1" sx={{ flex: 1 }}>
                              {member.ID}
                            </Typography>
                            <Chip 
                              label={`Categoría ${member.Category}`} 
                              size="small" 
                              color="primary"
                              variant="outlined"
                            />
                          </Stack>
                        </CardContent>
                      </Card>
                    ))}
                  </Stack>
                ) : (
                  <Typography variant="body2" color="text.secondary" sx={{ fontStyle: 'italic' }}>
                    Este equipo no tiene miembros registrados
                  </Typography>
                )}
              </Box>
                  </Stack>
                </Card>
              ))}
            </Stack>
          </Box>
        )}
      </CardContent>
    </Card>
  )
}

