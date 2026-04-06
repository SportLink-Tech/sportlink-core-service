import { useState, useEffect } from 'react'
import {
  Paper,
  Typography,
  Stack,
  FormControl,
  Select,
  MenuItem,
  TextField,
  Button,
  Box,
  Switch,
  Slider,
  Alert,
  CircularProgress,
  Tooltip,
} from '@mui/material'
import SearchIcon from '@mui/icons-material/Search'
import MyLocationIcon from '@mui/icons-material/MyLocation'
import { useGeolocation } from '../../../../shared/hooks/useGeolocation'
import { GeoFilter } from '../../../../shared/types/matchOffer'

interface MatchOfferFiltersProps {
  selectedSports: string[]
  fromDate: string
  toDate: string
  onFiltersChange: (filters: {
    sports: string[]
    fromDate: string
    toDate: string
    geoFilter?: GeoFilter
  }) => void
}

const AVAILABLE_SPORTS = ['Football', 'Paddle', 'Tennis'] as const

export function MatchOfferFilters({
  selectedSports,
  fromDate,
  toDate,
  onFiltersChange,
}: MatchOfferFiltersProps) {
  const [localSports, setLocalSports] = useState<string[]>(selectedSports)
  const [localFromDate, setLocalFromDate] = useState(fromDate)
  const [localToDate, setLocalToDate] = useState(toDate)
  const [geoEnabled, setGeoEnabled] = useState(false)
  const [radiusKm, setRadiusKm] = useState(50)
  const geo = useGeolocation()

  const handleGeoToggle = (enabled: boolean) => {
    setGeoEnabled(enabled)
    if (enabled && geo.status === 'idle') {
      geo.requestLocation()
    }
    if (!enabled) {
      geo.reset()
    }
  }

  const handleSearch = () => {
    const geoFilter: GeoFilter | undefined =
      geoEnabled && geo.status === 'granted' && geo.latitude !== null && geo.longitude !== null
        ? { latitude: geo.latitude, longitude: geo.longitude, radiusKm }
        : undefined

    onFiltersChange({
      sports: localSports,
      fromDate: localFromDate,
      toDate: localToDate,
      geoFilter,
    })
  }

  // Si el usuario habilitó geo pero luego fue denied/unavailable, desactivar el toggle
  useEffect(() => {
    if (geoEnabled && (geo.status === 'denied' || geo.status === 'unavailable')) {
      setGeoEnabled(false)
    }
  }, [geo.status, geoEnabled])

  // Sincronizar valores locales cuando cambian desde fuera
  useEffect(() => {
    setLocalSports(selectedSports)
    setLocalFromDate(fromDate)
    setLocalToDate(toDate)
  }, [selectedSports, fromDate, toDate])

  return (
    <Paper
      elevation={3}
      sx={{
        borderRadius: 4,
        overflow: 'hidden',
        boxShadow: '0 4px 20px rgba(0,0,0,0.1)',
      }}
    >
      <Box
        sx={{
          display: 'flex',
          flexDirection: { xs: 'column', md: 'row' },
          backgroundColor: 'white',
          alignItems: { md: 'flex-end' },
        }}
      >
        {/* Deporte */}
        <Box
          sx={{
            px: { xs: 3, md: 4 },
            py: { xs: 2.5, md: 3 },
            borderBottom: { xs: '1px solid #e0e0e0', md: 'none' },
            borderRight: { md: '1px solid #e0e0e0' },
            flex: { md: 1 },
            '&:hover': {
              backgroundColor: 'rgba(0, 0, 0, 0.02)',
            },
            transition: 'background-color 0.2s',
          }}
        >
          <Typography
            variant="caption"
            sx={{
              fontWeight: 700,
              color: 'text.secondary',
              textTransform: 'uppercase',
              fontSize: '0.7rem',
              letterSpacing: '0.8px',
              mb: 1,
              display: 'block',
            }}
          >
            Deporte
          </Typography>
          <FormControl fullWidth variant="standard">
            <Select
              multiple
              value={localSports}
              onChange={(e) => setLocalSports(e.target.value as string[])}
              displayEmpty
              sx={{
                '&:before': { display: 'none' },
                '&:after': { display: 'none' },
                '& .MuiSelect-select': {
                  py: 0.5,
                  px: 0,
                  fontSize: '1rem',
                  fontWeight: 500,
                  color: localSports.length === 0 ? '#9e9e9e' : 'text.primary',
                },
              }}
              renderValue={(selected) => {
                if (selected.length === 0) {
                  return <span style={{ color: '#9e9e9e' }}>Todos los deportes</span>
                }
                return selected.join(', ')
              }}
            >
              {AVAILABLE_SPORTS.map((sport) => (
                <MenuItem key={sport} value={sport}>
                  {sport}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Box>

        {/* Fechas */}
        <Box
          sx={{
            px: { xs: 3, md: 4 },
            py: { xs: 2.5, md: 3 },
            borderBottom: { xs: '1px solid #e0e0e0', md: 'none' },
            borderRight: { md: '1px solid #e0e0e0' },
            flex: { md: 1 },
          }}
        >
          <Typography
            variant="caption"
            sx={{
              fontWeight: 700,
              color: 'text.secondary',
              textTransform: 'uppercase',
              fontSize: '0.7rem',
              letterSpacing: '0.8px',
              mb: 1,
              display: 'block',
            }}
          >
            Fechas
          </Typography>
          <Stack direction={{ xs: 'column', md: 'row' }} spacing={{ xs: 1.5, md: 2 }}>
            <TextField
              type="date"
              value={localFromDate}
              onChange={(e) => setLocalFromDate(e.target.value)}
              placeholder="Desde"
              variant="standard"
              fullWidth
              InputProps={{
                disableUnderline: true,
              }}
              sx={{
                '& .MuiInputBase-input': {
                  py: 0.5,
                  px: 0,
                  fontSize: '1rem',
                  fontWeight: 500,
                  color: localFromDate ? 'text.primary' : '#9e9e9e',
                },
              }}
            />
            <TextField
              type="date"
              value={localToDate}
              onChange={(e) => setLocalToDate(e.target.value)}
              placeholder="Hasta"
              variant="standard"
              fullWidth
              InputProps={{
                disableUnderline: true,
              }}
              sx={{
                '& .MuiInputBase-input': {
                  py: 0.5,
                  px: 0,
                  fontSize: '1rem',
                  fontWeight: 500,
                  color: localToDate ? 'text.primary' : '#9e9e9e',
                },
              }}
            />
          </Stack>
        </Box>

        {/* Cerca de mí */}
        <Box
          sx={{
            px: { xs: 3, md: 4 },
            py: { xs: 2.5, md: 3 },
            borderBottom: { xs: '1px solid #e0e0e0', md: 'none' },
            borderRight: { md: '1px solid #e0e0e0' },
            flex: { md: 1 },
          }}
        >
          <Stack direction="row" alignItems="center" justifyContent="space-between">
            <Typography
              variant="caption"
              sx={{
                fontWeight: 700,
                color: 'text.secondary',
                textTransform: 'uppercase',
                fontSize: '0.7rem',
                letterSpacing: '0.8px',
              }}
            >
              Cerca de mí
            </Typography>
            <Tooltip title={geo.status === 'unavailable' ? 'Tu dispositivo no soporta geolocalización' : ''}>
              <span>
                <Switch
                  size="small"
                  checked={geoEnabled}
                  onChange={(e) => handleGeoToggle(e.target.checked)}
                  disabled={geo.status === 'unavailable' || geo.status === 'loading'}
                />
              </span>
            </Tooltip>
          </Stack>

          {geoEnabled && geo.status === 'loading' && (
            <Stack direction="row" spacing={1} alignItems="center" sx={{ mt: 1 }}>
              <CircularProgress size={14} />
              <Typography variant="caption" color="text.secondary">Detectando ubicación...</Typography>
            </Stack>
          )}

          {geoEnabled && geo.status === 'granted' && (
            <Box sx={{ mt: 1, px: 0.5 }}>
              <Stack direction="row" alignItems="center" spacing={1}>
                <MyLocationIcon sx={{ fontSize: 14, color: 'success.main' }} />
                <Typography variant="caption" color="text.secondary" sx={{ flex: 1 }}>
                  Radio: <strong>{radiusKm} km</strong>
                </Typography>
              </Stack>
              <Slider
                value={radiusKm}
                onChange={(_, v) => setRadiusKm(v as number)}
                min={5}
                max={100}
                step={5}
                size="small"
                sx={{ mt: 0.5 }}
              />
            </Box>
          )}

          {geo.status === 'denied' && (
            <Alert severity="warning" sx={{ mt: 1, py: 0 }} icon={false}>
              <Typography variant="caption">
                Permiso denegado. Habilitalo en tu navegador.
              </Typography>
            </Alert>
          )}
        </Box>

        {/* Botón Buscar - Grande y destacado */}
        <Box
          sx={{
            p: { xs: 2.5, md: 2 },
            backgroundColor: 'white',
            display: 'flex',
            alignItems: { md: 'flex-end' },
          }}
        >
          <Button
            variant="contained"
            size="large"
            onClick={handleSearch}
            startIcon={<SearchIcon sx={{ fontSize: 24 }} />}
            sx={{
              borderRadius: { xs: 3, md: '50%' },
              minWidth: { xs: '100%', md: 64 },
              width: { xs: '100%', md: 64 },
              height: { xs: 56, md: 64 },
              fontSize: { xs: '1.1rem', md: '1.5rem' },
              fontWeight: 700,
              textTransform: 'none',
              background: 'linear-gradient(135deg, #00C853 0%, #6A1B9A 100%)',
              boxShadow: '0 4px 15px rgba(0, 200, 83, 0.3)',
              '&:hover': {
                background: 'linear-gradient(135deg, #009624 0%, #38006B 100%)',
                transform: 'translateY(-2px)',
                boxShadow: '0 6px 20px rgba(0, 200, 83, 0.4)',
              },
              transition: 'all 0.2s ease-in-out',
              '& .MuiButton-startIcon': {
                margin: { xs: '0 8px 0 0', md: 0 },
              },
              '& .MuiButton-startIcon > *:nth-of-type(1)': {
                fontSize: { xs: 24, md: 28 },
              },
            }}
          >
            <Box component="span" sx={{ display: { xs: 'inline', md: 'none' } }}>
              Buscar
            </Box>
          </Button>
        </Box>
      </Box>
    </Paper>
  )
}

