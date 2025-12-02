import { useState } from 'react'
import {
  Box,
  Button,
  Card,
  CardContent,
  TextField,
  Typography,
  Alert,
  Chip,
  Stack,
  CircularProgress,
  Divider,
  Paper,
  FormGroup,
  FormControlLabel,
  Checkbox,
  FormLabel,
} from '@mui/material'
import SearchIcon from '@mui/icons-material/Search'
import { SportSelect } from '../../../../shared/components/atoms/SportSelect'
import { Sport, Team } from '../../../../shared/types/team'
import { useTeamContext } from '../../context/TeamContext'

const SPORTS: Sport[] = ['Football', 'Paddle', 'Tennis']

/**
 * Feature Page: Search Team
 * Uses SearchTeamUseCase from context
 * Following Atomic Hexagonal Architecture
 */
export function SearchTeamPage() {
  const { searchTeamUseCase } = useTeamContext()

  const [sport, setSport] = useState<Sport>('Paddle')
  const [teamName, setTeamName] = useState('')
  const [selectedCategories, setSelectedCategories] = useState<number[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [teams, setTeams] = useState<Team[]>([])

  const categories = [1, 2, 3, 4, 5, 6, 7]

  const handleCategoryToggle = (category: number) => {
    setSelectedCategories(prev =>
      prev.includes(category)
        ? prev.filter(c => c !== category)
        : [...prev, category]
    )
  }

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setTeams([])
    setLoading(true)

    // Execute use case
    const result = await searchTeamUseCase.execute(sport, teamName || undefined, selectedCategories.length > 0 ? selectedCategories : undefined)

    if (result.success) {
      setTeams(result.teams)
      if (result.teams.length === 0) {
        setError('No se encontraron equipos con los criterios especificados')
      }
    } else {
      setError(result.error || 'Error al buscar equipos')
    }

    setLoading(false)
  }

  return (
    <Box>
      <Stack spacing={4}>
        {/* Hero Section */}
        <Paper
          elevation={0}
          sx={{
            p: 4,
            background: 'linear-gradient(135deg, #00C853 0%, #6A1B9A 100%)',
            color: 'white',
            borderRadius: 4,
          }}
        >
          <Stack spacing={2} alignItems="center">
            <SearchIcon sx={{ fontSize: 64 }} />
            <Typography variant="h3" component="h1" align="center" fontWeight={700}>
              Encuentra tu Equipo
            </Typography>
            <Typography variant="h6" align="center" sx={{ opacity: 0.95, maxWidth: 600 }}>
              Busca equipos por deporte, nombre y categoría. Puedes combinar los filtros o buscar solo por deporte para ver todos los equipos.
            </Typography>
          </Stack>
        </Paper>

        {/* Search Card */}
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
                <SportSelect
                  value={sport}
                  onChange={setSport}
                  required
                  sports={SPORTS}
                />

                <TextField
                  label="Nombre del Equipo (opcional)"
                  value={teamName}
                  onChange={(e) => setTeamName(e.target.value)}
                  fullWidth
                  placeholder="Ej: Thunder Strikers"
                  helperText="Deja vacío para buscar todos los equipos del deporte"
                />

                <Box>
                  <FormLabel component="legend" sx={{ mb: 1 }}>
                    Categorías (opcional)
                  </FormLabel>
                  <FormGroup row>
                    {categories.map((category) => (
                      <FormControlLabel
                        key={category}
                        control={
                          <Checkbox
                            checked={selectedCategories.includes(category)}
                            onChange={() => handleCategoryToggle(category)}
                            color="primary"
                          />
                        }
                        label={`L${category}`}
                      />
                    ))}
                  </FormGroup>
                  {selectedCategories.length > 0 && (
                    <Box sx={{ mt: 1 }}>
                      <Typography variant="caption" color="text.secondary">
                        Seleccionadas: {selectedCategories.map(c => `L${c}`).join(', ')}
                      </Typography>
                    </Box>
                  )}
                </Box>

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
      </Stack>
    </Box>
  )
}

