import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Box,
  Button,
  Card,
  CardContent,
  TextField,
  Typography,
  Stack,
  CircularProgress,
  Snackbar,
  Alert,
  Paper,
  MenuItem,
  Chip,
  FormControl,
  InputLabel,
  Select,
  SelectChangeEvent,
  Grid,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Autocomplete,
  Grow,
  Fade,
} from '@mui/material'
import AddCircleIcon from '@mui/icons-material/AddCircle'
import EventIcon from '@mui/icons-material/Event'
import CheckCircleIcon from '@mui/icons-material/CheckCircle'
import LocationOnIcon from '@mui/icons-material/LocationOn'
import AccessTimeIcon from '@mui/icons-material/AccessTime'
import GroupsIcon from '@mui/icons-material/Groups'
import { SportSelect } from '../../../../shared/components/atoms/SportSelect'
import { Sport } from '../../../../shared/types/team'
import { useMatchAnnouncementContext } from '../../context/MatchAnnouncementContext'
import { MatchAnnouncement } from '../../../../shared/types/matchAnnouncement'

const SPORTS: Sport[] = ['Football', 'Paddle', 'Tennis']

export function CreateMatchAnnouncementPage() {
  const { createMatchAnnouncementUseCase } = useMatchAnnouncementContext()
  const navigate = useNavigate()

  const [sport, setSport] = useState<Sport>('Paddle')
  const [teamName, setTeamName] = useState('')
  const [availableTeams, setAvailableTeams] = useState<string[]>([])
  const [loadingTeams, setLoadingTeams] = useState(false)
  const [day, setDay] = useState('')
  const [startTime, setStartTime] = useState('')
  const [endTime, setEndTime] = useState('')
  const [country, setCountry] = useState('Argentina')
  const [province, setProvince] = useState('Buenos Aires')
  const [locality, setLocality] = useState('')
  
  const [categoryRangeType, setCategoryRangeType] = useState<'SPECIFIC' | 'GREATER_THAN' | 'LESS_THAN' | 'BETWEEN'>('SPECIFIC')
  const [selectedCategories, setSelectedCategories] = useState<number[]>([])
  const [minLevel, setMinLevel] = useState<number>(1)
  const [maxLevel, setMaxLevel] = useState<number>(7)

  const [loading, setLoading] = useState(false)
  const [showSuccessDialog, setShowSuccessDialog] = useState(false)
  const [createdAnnouncement, setCreatedAnnouncement] = useState<MatchAnnouncement | null>(null)
  const [showErrorSnackbar, setShowErrorSnackbar] = useState(false)
  const [errorMessage, setErrorMessage] = useState('')
  const [attempted, setAttempted] = useState(false)

  // Load teams when sport changes
  useEffect(() => {
    const fetchTeams = async () => {
      setLoadingTeams(true)
      try {
        const response = await fetch(`/sport/${sport}/team`)
        if (response.ok) {
          const teams = await response.json()
          const teamNames = teams.map((team: any) => team.Name)
          setAvailableTeams(teamNames)
        } else {
          setAvailableTeams([])
        }
      } catch (error) {
        console.error('Error fetching teams:', error)
        setAvailableTeams([])
      } finally {
        setLoadingTeams(false)
      }
    }

    fetchTeams()
  }, [sport])

  const handleCategoryChange = (category: number) => {
    setSelectedCategories((prev) =>
      prev.includes(category) ? prev.filter((c) => c !== category) : [...prev, category].sort((a, b) => a - b)
    )
  }

  const handleCategoryRangeTypeChange = (event: SelectChangeEvent) => {
    setCategoryRangeType(event.target.value as any)
    setSelectedCategories([])
    setMinLevel(1)
    setMaxLevel(7)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setAttempted(true)
    
    if (!isFormValid()) {
      setErrorMessage('Por favor completa todos los campos obligatorios')
      setShowErrorSnackbar(true)
      return
    }
    
    setLoading(true)
    setErrorMessage('')

    // Build admitted categories (snake_case)
    const admittedCategories: any = {
      type: categoryRangeType,
    }

    if (categoryRangeType === 'SPECIFIC') {
      admittedCategories.categories = selectedCategories
    } else if (categoryRangeType === 'GREATER_THAN') {
      admittedCategories.min_level = minLevel
    } else if (categoryRangeType === 'LESS_THAN') {
      admittedCategories.max_level = maxLevel
    } else if (categoryRangeType === 'BETWEEN') {
      admittedCategories.min_level = minLevel
      admittedCategories.max_level = maxLevel
    }

    const result = await createMatchAnnouncementUseCase.execute({
      team_name: teamName,
      sport,
      day,
      time_slot: {
        start_time: `${day}T${startTime}:00`,
        end_time: `${day}T${endTime}:00`,
      },
      location: {
        country,
        province,
        locality,
      },
      admitted_categories: admittedCategories,
    })

    if (result.success) {
      setCreatedAnnouncement(result.announcement)
      setShowSuccessDialog(true)
      // Reset form
      setTeamName('')
      setDay('')
      setStartTime('')
      setEndTime('')
      setLocality('')
      setSelectedCategories([])
      setAttempted(false)
    } else {
      setErrorMessage(result.error || 'Error al crear el anuncio')
      setShowErrorSnackbar(true)
    }

    setLoading(false)
  }

  const handleSuccessDialogClose = () => {
    setShowSuccessDialog(false)
    setCreatedAnnouncement(null)
    navigate('/')
  }

  const getCategoryText = (admittedCategories: MatchAnnouncement['admitted_categories']) => {
    if (!admittedCategories) {
      return 'No especificado'
    }

    switch (admittedCategories.type) {
      case 'SPECIFIC':
        return admittedCategories.categories && admittedCategories.categories.length > 0
          ? `Categorías: ${admittedCategories.categories.map(c => `L${c}`).join(', ')}`
          : 'Categorías no especificadas'
      case 'GREATER_THAN':
        return admittedCategories.min_level
          ? `Nivel >= L${admittedCategories.min_level}`
          : 'Nivel mínimo no especificado'
      case 'LESS_THAN':
        return admittedCategories.max_level
          ? `Nivel <= L${admittedCategories.max_level}`
          : 'Nivel máximo no especificado'
      case 'BETWEEN':
        return admittedCategories.min_level && admittedCategories.max_level
          ? `Nivel L${admittedCategories.min_level} - L${admittedCategories.max_level}`
          : 'Rango no especificado'
      default:
        return 'Cualquier nivel'
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('es-AR', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
  }

  const formatTime = (dateTimeString: string) => {
    const date = new Date(dateTimeString)
    return date.toLocaleTimeString('es-AR', { hour: '2-digit', minute: '2-digit' })
  }

  const isFormValid = () => {
    return teamName && day && startTime && endTime && locality &&
      (categoryRangeType === 'SPECIFIC' ? selectedCategories.length > 0 : true)
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
                Publicar Partido
              </Typography>
              <Typography variant="h6" align="center" sx={{ opacity: 0.95, maxWidth: 600 }}>
                Anuncia tu intención de jugar un partido y encuentra rivales
              </Typography>
            </Stack>
          </Paper>

          {/* Form Card */}
          <Card>
            <CardContent>
              <form onSubmit={handleSubmit}>
                <Stack spacing={3}>
                  {/* Team Name - Autocomplete */}
                  <Autocomplete
                    value={teamName}
                    onChange={(_, newValue) => setTeamName(newValue || '')}
                    inputValue={teamName}
                    onInputChange={(_, newInputValue) => setTeamName(newInputValue)}
                    options={availableTeams}
                    loading={loadingTeams}
                    freeSolo
                    renderInput={(params) => (
                      <TextField
                        {...params}
                        label="Nombre del Equipo"
                        required
                        placeholder="Ej: Boca Junior"
                        error={attempted && !teamName}
                        helperText={
                          attempted && !teamName
                            ? "Campo obligatorio"
                            : attempted && availableTeams.length > 0 && !availableTeams.includes(teamName)
                            ? "⚠️ Este equipo no existe. Verifica el nombre o créalo primero."
                            : ""
                        }
                        InputProps={{
                          ...params.InputProps,
                          endAdornment: (
                            <>
                              {loadingTeams ? <CircularProgress color="inherit" size={20} /> : null}
                              {params.InputProps.endAdornment}
                            </>
                          ),
                        }}
                      />
                    )}
                  />

                  {/* Sport */}
                  <SportSelect sports={SPORTS} value={sport} onChange={setSport} />

                  {/* Date */}
                  <TextField
                    label="Día del Partido"
                    type="date"
                    value={day}
                    onChange={(e) => setDay(e.target.value)}
                    fullWidth
                    required
                    InputLabelProps={{ shrink: true }}
                    inputProps={{ min: new Date().toISOString().split('T')[0] }}
                    error={attempted && !day}
                    helperText={attempted && !day ? "Campo obligatorio" : ""}
                  />

                  {/* Time Slot */}
                  <Grid container spacing={2}>
                    <Grid item xs={6}>
                      <TextField
                        label="Hora Inicio"
                        type="time"
                        value={startTime}
                        onChange={(e) => setStartTime(e.target.value)}
                        fullWidth
                        required
                        InputLabelProps={{ shrink: true }}
                        error={attempted && !startTime}
                        helperText={attempted && !startTime ? "Campo obligatorio" : ""}
                      />
                    </Grid>
                    <Grid item xs={6}>
                      <TextField
                        label="Hora Fin"
                        type="time"
                        value={endTime}
                        onChange={(e) => setEndTime(e.target.value)}
                        fullWidth
                        required
                        InputLabelProps={{ shrink: true }}
                        error={attempted && !endTime}
                        helperText={attempted && !endTime ? "Campo obligatorio" : ""}
                      />
                    </Grid>
                  </Grid>

                  {/* Location */}
                  <Typography variant="subtitle1" fontWeight={600} sx={{ mt: 2 }}>
                    Ubicación
                  </Typography>
                  <Grid container spacing={2}>
                    <Grid item xs={4}>
                      <TextField
                        label="País"
                        value={country}
                        onChange={(e) => setCountry(e.target.value)}
                        fullWidth
                        required
                      />
                    </Grid>
                    <Grid item xs={4}>
                      <TextField
                        label="Provincia"
                        value={province}
                        onChange={(e) => setProvince(e.target.value)}
                        fullWidth
                        required
                      />
                    </Grid>
                    <Grid item xs={4}>
                      <TextField
                        label="Localidad"
                        value={locality}
                        onChange={(e) => setLocality(e.target.value)}
                        fullWidth
                        required
                        placeholder="Ej: CABA"
                        error={attempted && !locality}
                        helperText={attempted && !locality ? "Campo obligatorio" : ""}
                      />
                    </Grid>
                  </Grid>

                  {/* Category Range */}
                  <Typography variant="subtitle1" fontWeight={600} sx={{ mt: 2 }}>
                    Categorías Admitidas
                  </Typography>
                  
                  <FormControl fullWidth>
                    <InputLabel>Tipo de Rango</InputLabel>
                    <Select value={categoryRangeType} onChange={handleCategoryRangeTypeChange} label="Tipo de Rango">
                      <MenuItem value="SPECIFIC">Categorías Específicas</MenuItem>
                      <MenuItem value="GREATER_THAN">Mayor o igual a</MenuItem>
                      <MenuItem value="LESS_THAN">Menor o igual a</MenuItem>
                      <MenuItem value="BETWEEN">Entre dos niveles</MenuItem>
                    </Select>
                  </FormControl>

                  {categoryRangeType === 'SPECIFIC' && (
                    <>
                      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                        {[1, 2, 3, 4, 5, 6, 7].map((cat) => (
                          <Chip
                            key={cat}
                            label={`L${cat}`}
                            color={selectedCategories.includes(cat) ? 'primary' : 'default'}
                            onClick={() => handleCategoryChange(cat)}
                            variant={selectedCategories.includes(cat) ? 'filled' : 'outlined'}
                          />
                        ))}
                      </Box>
                      {selectedCategories.length > 0 ? (
                        <Typography variant="body2" color="text.secondary">
                          Seleccionadas: {selectedCategories.map((c) => `L${c}`).join(', ')}
                        </Typography>
                      ) : attempted ? (
                        <Typography variant="body2" color="error">
                          Debes seleccionar al menos una categoría
                        </Typography>
                      ) : null}
                    </>
                  )}

                  {categoryRangeType === 'GREATER_THAN' && (
                    <TextField
                      label="Nivel Mínimo"
                      type="number"
                      value={minLevel}
                      onChange={(e) => setMinLevel(Number(e.target.value))}
                      inputProps={{ min: 1, max: 7 }}
                      fullWidth
                      helperText="Ej: 5 significa L5 o superior"
                    />
                  )}

                  {categoryRangeType === 'LESS_THAN' && (
                    <TextField
                      label="Nivel Máximo"
                      type="number"
                      value={maxLevel}
                      onChange={(e) => setMaxLevel(Number(e.target.value))}
                      inputProps={{ min: 1, max: 7 }}
                      fullWidth
                      helperText="Ej: 3 significa L3 o inferior"
                    />
                  )}

                  {categoryRangeType === 'BETWEEN' && (
                    <Grid container spacing={2}>
                      <Grid item xs={6}>
                        <TextField
                          label="Nivel Mínimo"
                          type="number"
                          value={minLevel}
                          onChange={(e) => setMinLevel(Number(e.target.value))}
                          inputProps={{ min: 1, max: 7 }}
                          fullWidth
                        />
                      </Grid>
                      <Grid item xs={6}>
                        <TextField
                          label="Nivel Máximo"
                          type="number"
                          value={maxLevel}
                          onChange={(e) => setMaxLevel(Number(e.target.value))}
                          inputProps={{ min: 1, max: 7 }}
                          fullWidth
                        />
                      </Grid>
                    </Grid>
                  )}

                  {/* Submit Button */}
                  <Button
                    type="submit"
                    variant="contained"
                    size="large"
                    disabled={loading || !isFormValid()}
                    startIcon={loading ? <CircularProgress size={20} /> : <EventIcon />}
                    sx={{ mt: 2 }}
                  >
                    {loading ? 'Publicando...' : 'Publicar Partido'}
                  </Button>
                </Stack>
              </form>
            </CardContent>
          </Card>
        </Stack>
      </Box>

      {/* Success Dialog - Minimalist */}
      <Dialog 
        open={showSuccessDialog} 
        onClose={handleSuccessDialogClose}
        maxWidth="sm"
        fullWidth
        TransitionComponent={Grow}
        transitionDuration={500}
      >
        <DialogTitle sx={{ bgcolor: 'success.main', color: 'white', textAlign: 'center', py: 2 }}>
          ¡Partido Publicado!
        </DialogTitle>
        <DialogContent sx={{ pt: 3, pb: 2 }}>
          {createdAnnouncement && (
            <Fade in timeout={800}>
              <Stack spacing={2.5}>
                {/* Date and Time */}
                <Box>
                  <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 0.5 }}>
                    <EventIcon fontSize="small" color="action" />
                    <Typography variant="caption" color="text.secondary" textTransform="uppercase">
                      Fecha y Horario
                    </Typography>
                  </Stack>
                  <Typography variant="body1">
                    {formatDate(createdAnnouncement.day)}
                  </Typography>
                  <Stack direction="row" spacing={1} alignItems="center" sx={{ mt: 0.5 }}>
                    <AccessTimeIcon fontSize="small" color="action" />
                    <Typography variant="body2">
                      {formatTime(createdAnnouncement.time_slot.start_time)} - {formatTime(createdAnnouncement.time_slot.end_time)}
                    </Typography>
                  </Stack>
                </Box>

                {/* Location */}
                <Box>
                  <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 0.5 }}>
                    <LocationOnIcon fontSize="small" color="action" />
                    <Typography variant="caption" color="text.secondary" textTransform="uppercase">
                      Ubicación
                    </Typography>
                  </Stack>
                  <Typography variant="body1">
                    {createdAnnouncement.location.locality}, {createdAnnouncement.location.province}
                  </Typography>
                </Box>

                {/* Categories */}
                <Box>
                  <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 0.5 }}>
                    <GroupsIcon fontSize="small" color="action" />
                    <Typography variant="caption" color="text.secondary" textTransform="uppercase">
                      Categorías
                    </Typography>
                  </Stack>
                  <Typography variant="body1">
                    {getCategoryText(createdAnnouncement.admitted_categories)}
                  </Typography>
                </Box>
              </Stack>
            </Fade>
          )}
        </DialogContent>
        <DialogActions sx={{ p: 2, pt: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2 }}>
          <Grow in timeout={1000}>
            <CheckCircleIcon sx={{ fontSize: 52, color: 'success.main' }} />
          </Grow>
          <Button 
            onClick={handleSuccessDialogClose} 
            variant="contained" 
            size="large"
            fullWidth
          >
            Ver Partidos Disponibles
          </Button>
        </DialogActions>
      </Dialog>

      {/* Error Snackbar */}
      <Snackbar
        open={showErrorSnackbar}
        autoHideDuration={6000}
        onClose={() => setShowErrorSnackbar(false)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert onClose={() => setShowErrorSnackbar(false)} severity="error" sx={{ width: '100%' }}>
          {errorMessage}
        </Alert>
      </Snackbar>
    </>
  )
}

