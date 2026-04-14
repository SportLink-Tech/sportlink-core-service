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
  Grow,
  Fade,
} from '@mui/material'
import AddCircleIcon from '@mui/icons-material/AddCircle'
import ArrowBackIcon from '@mui/icons-material/ArrowBack'
import EventIcon from '@mui/icons-material/Event'
import CheckCircleIcon from '@mui/icons-material/CheckCircle'
import LocationOnIcon from '@mui/icons-material/LocationOn'
import AccessTimeIcon from '@mui/icons-material/AccessTime'
import GroupsIcon from '@mui/icons-material/Groups'
import MyLocationIcon from '@mui/icons-material/MyLocation'
import { SportSelect } from '../../../../shared/components/atoms/SportSelect'
import { Sport, Team } from '../../../../shared/types/team'
import { useMatchOfferContext } from '../../context/MatchOfferContext'
import { MatchOffer } from '../../../../shared/types/matchOffer'
import { useGeolocation } from '../../../../shared/hooks/useGeolocation'
import { useTeamContext } from '../../../team/context/TeamContext'
import { useAuth } from '../../../auth/context/AuthContext'

const SPORTS: Sport[] = ['Football', 'Paddle', 'Tennis']

export function CreateMatchOfferPage() {
  const { createMatchOfferUseCase } = useMatchOfferContext()
  const { listAccountTeamsUseCase } = useTeamContext()
  const { accountId } = useAuth()
  const navigate = useNavigate()

  const [sport, setSport] = useState<Sport>('Paddle')
  const handleSportChange = (newSport: Sport) => { setSport(newSport); setTeamName('') }
  const [teamName, setTeamName] = useState('')
  const [userTeams, setUserTeams] = useState<Team[]>([])
  const [loadingTeams, setLoadingTeams] = useState(false)
  const today = new Date()
  const todayStr = today.toISOString().split('T')[0]
  const oneHourAhead = new Date(today.getTime() + 60 * 60 * 1000)
  const startTimeDefault = `${String(oneHourAhead.getHours()).padStart(2, '0')}:${String(oneHourAhead.getMinutes()).padStart(2, '0')}`

  const [day, setDay] = useState(todayStr)
  const [startTime, setStartTime] = useState(startTimeDefault)
  const [endTime, setEndTime] = useState('23:59')
  const [country, setCountry] = useState('Argentina')
  const [province, setProvince] = useState('Buenos Aires')
  const [locality, setLocality] = useState('')
  const geo = useGeolocation()
  
  const [capacity, setCapacity] = useState<number>(0)
  const [categoryRangeType, setCategoryRangeType] = useState<'SPECIFIC' | 'GREATER_THAN' | 'LESS_THAN' | 'BETWEEN'>('SPECIFIC')
  const [selectedCategories, setSelectedCategories] = useState<number[]>([])
  const [minLevel, setMinLevel] = useState<number>(1)
  const [maxLevel, setMaxLevel] = useState<number>(7)

  const [reverseGeoLoading, setReverseGeoLoading] = useState(false)
  const [loading, setLoading] = useState(false)
  const [showSuccessDialog, setShowSuccessDialog] = useState(false)
  const [createdAnnouncement, setCreatedAnnouncement] = useState<MatchOffer | null>(null)
  const [showErrorSnackbar, setShowErrorSnackbar] = useState(false)
  const [errorMessage, setErrorMessage] = useState('')
  const [attempted, setAttempted] = useState(false)

  // Reverse geocode when geolocation is granted
  useEffect(() => {
    if (geo.status === 'granted' && geo.latitude !== null && geo.longitude !== null) {
      setReverseGeoLoading(true)
      fetch(
        `https://nominatim.openstreetmap.org/reverse?lat=${geo.latitude}&lon=${geo.longitude}&format=json`,
        { headers: { 'Accept-Language': 'es' } }
      )
        .then((res) => res.json())
        .then((data) => {
          const addr = data.address
          if (addr?.country) setCountry(addr.country)
          if (addr?.state) setProvince(addr.state)
          const loc = addr?.city || addr?.town || addr?.village || addr?.municipality || addr?.suburb || ''
          if (loc) setLocality(loc)
        })
        .catch(console.error)
        .finally(() => setReverseGeoLoading(false))
    }
  }, [geo.status, geo.latitude, geo.longitude])

  // Load user teams on mount
  useEffect(() => {
    const fetchTeams = async () => {
      setLoadingTeams(true)
      const result = await listAccountTeamsUseCase.execute(accountId ?? '')
      setUserTeams(result.teams)
      setLoadingTeams(false)
    }
    fetchTeams()
  }, [])

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

    const result = await createMatchOfferUseCase.execute(accountId ?? '', {
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
        ...(geo.status === 'granted' && geo.latitude !== null && geo.longitude !== null
          ? { latitude: geo.latitude, longitude: geo.longitude }
          : {}),
      },
      admitted_categories: admittedCategories,
      capacity,
    })

    if (result.success) {
      setCreatedAnnouncement(result.announcement)
      setShowSuccessDialog(true)
      // Reset form
      setTeamName('')
      setDay(new Date().toISOString().split('T')[0])
      const nextHour = new Date(Date.now() + 60 * 60 * 1000)
      setStartTime(`${String(nextHour.getHours()).padStart(2, '0')}:${String(nextHour.getMinutes()).padStart(2, '0')}`)
      setEndTime('23:59')
      setLocality('')
      setCapacity(0)
      setSelectedCategories([])
      setAttempted(false)
    } else {
      setErrorMessage(result.error || 'Error al crear la oferta')
      setShowErrorSnackbar(true)
    }

    setLoading(false)
  }

  const handleSuccessDialogClose = () => {
    setShowSuccessDialog(false)
    setCreatedAnnouncement(null)
    navigate('/my-offers')
  }

  const getCategoryText = (admittedCategories: MatchOffer['admitted_categories']) => {
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

  const isTimeSlotValid = () => {
    if (!startTime || !endTime) return true // campo vacío, lo maneja otro check
    return endTime > startTime
  }

  const isFormValid = () => {
    return day && startTime && endTime && locality && isTimeSlotValid() &&
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
              <Button
                variant="contained"
                startIcon={<ArrowBackIcon />}
                onClick={() => navigate('/my-offers')}
                sx={{ bgcolor: 'white', color: 'primary.main', '&:hover': { bgcolor: 'grey.100' } }}
              >
                Volver a mis publicaciones
              </Button>
            </Stack>
          </Paper>

          {/* Form Card */}
          <Card>
            <CardContent>
              <form onSubmit={handleSubmit}>
                <Stack spacing={3}>
                  {/* Team Name - Optional select from user's teams */}
                  <FormControl fullWidth>
                    <InputLabel>Equipo (opcional)</InputLabel>
                    <Select
                      value={teamName}
                      onChange={(e) => setTeamName(e.target.value)}
                      label="Equipo (opcional)"
                      disabled={loadingTeams}
                      startAdornment={loadingTeams ? <CircularProgress size={20} sx={{ mr: 1 }} /> : null}
                    >
                      <MenuItem value="">Sin equipo</MenuItem>
                      {userTeams
                        .filter((t) => t.Sport === sport)
                        .map((t) => (
                          <MenuItem key={t.Name} value={t.Name}>
                            {t.Name}
                          </MenuItem>
                        ))}
                      {!loadingTeams && userTeams.filter((t) => t.Sport === sport).length === 0 && (
                        <MenuItem disabled value="">
                          No tenés equipos de {sport}
                        </MenuItem>
                      )}
                    </Select>
                  </FormControl>

                  {/* Sport */}
                  <SportSelect sports={SPORTS} value={sport} onChange={handleSportChange} />

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
                        inputProps={{ step: 60, lang: 'de' }}
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
                        inputProps={{ step: 60, lang: 'de' }}
                        error={attempted && (!endTime || !isTimeSlotValid())}
                        helperText={
                          attempted && !endTime
                            ? "Campo obligatorio"
                            : attempted && !isTimeSlotValid()
                            ? "La hora de fin debe ser posterior a la de inicio"
                            : ""
                        }
                      />
                    </Grid>
                  </Grid>

                  {/* Capacity */}
                  <TextField
                    label="Cupos del partido"
                    type="number"
                    value={capacity}
                    onChange={(e) => setCapacity(Math.max(0, Number(e.target.value)))}
                    fullWidth
                    InputLabelProps={{ shrink: true }}
                    inputProps={{ min: 0 }}
                    helperText={
                      capacity === 0
                        ? 'Sin límite de cupos ni confirmación automática'
                        : `${capacity} cupo${capacity !== 1 ? 's' : ''} en total (incluye al organizador). Al completarse se confirmará automáticamente.`
                    }
                    InputProps={{
                      startAdornment: <GroupsIcon sx={{ mr: 1, color: 'action.active' }} />,
                    }}
                  />

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

                  {/* Geolocation */}
                  <Box>
                    <Button
                      variant="outlined"
                      size="small"
                      startIcon={
                        geo.status === 'loading'
                          ? <CircularProgress size={16} />
                          : geo.status === 'granted'
                          ? <CheckCircleIcon color="success" />
                          : <MyLocationIcon />
                      }
                      onClick={geo.status === 'granted' ? geo.reset : geo.requestLocation}
                      disabled={geo.status === 'loading'}
                      color={geo.status === 'granted' ? 'success' : 'primary'}
                    >
                      {geo.status === 'loading' && 'Detectando ubicación...'}
                      {geo.status === 'granted' && (reverseGeoLoading ? 'Obteniendo dirección...' : 'Ubicación detectada')}
                      {(geo.status === 'idle' || geo.status === 'denied' || geo.status === 'unavailable') && 'Usar mi ubicación actual'}
                    </Button>

                    {geo.status === 'denied' && (
                      <Alert severity="warning" sx={{ mt: 1 }} icon={false}>
                        No se pudo acceder a tu ubicación. El partido se publicará sin coordenadas y no aparecerá en búsquedas por proximidad. Podés habilitarlo desde la configuración de tu navegador.
                      </Alert>
                    )}
                    {geo.status === 'unavailable' && (
                      <Alert severity="info" sx={{ mt: 1 }} icon={false}>
                        Tu dispositivo no soporta geolocalización. El partido se publicará sin coordenadas.
                      </Alert>
                    )}
                    {geo.status === 'idle' && (
                      <Typography variant="caption" color="text.secondary" display="block" sx={{ mt: 0.5 }}>
                        Opcional. Sin coordenadas el partido no aparece en búsquedas por proximidad.
                      </Typography>
                    )}
                  </Box>

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

