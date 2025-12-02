import { useState, useEffect } from 'react'
import {
  Box,
  Typography,
  Stack,
  Paper,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Alert,
  Grid,
  Divider,
} from '@mui/material'
import EventIcon from '@mui/icons-material/Event'
import LocationOnIcon from '@mui/icons-material/LocationOn'
import SportsIcon from '@mui/icons-material/Sports'
import GroupsIcon from '@mui/icons-material/Groups'
import AccessTimeIcon from '@mui/icons-material/AccessTime'
import { useMatchAnnouncementContext } from '../../context/MatchAnnouncementContext'
import { MatchAnnouncement } from '../../../../shared/types/matchAnnouncement'

export function ListMatchAnnouncementsPage() {
  const { findMatchAnnouncementsUseCase } = useMatchAnnouncementContext()

  const [announcements, setAnnouncements] = useState<MatchAnnouncement[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadAnnouncements()
  }, [])

  const loadAnnouncements = async () => {
    setLoading(true)
    setError(null)

    // La query por defecto incluye fromDate = hoy (se setea automáticamente en el use case)
    const result = await findMatchAnnouncementsUseCase.execute({})

    if (result.success) {
      setAnnouncements(result.announcements)
    } else {
      setError(result.error || 'Error al cargar los anuncios')
    }

    setLoading(false)
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('es-AR', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
  }

  const formatTime = (dateTimeString: string) => {
    const date = new Date(dateTimeString)
    return date.toLocaleTimeString('es-AR', { hour: '2-digit', minute: '2-digit' })
  }

  const getCategoryText = (admittedCategories: MatchAnnouncement['admitted_categories']) => {
    switch (admittedCategories.type) {
      case 'SPECIFIC':
        return `Categorías: ${admittedCategories.categories?.map(c => `L${c}`).join(', ')}`
      case 'GREATER_THAN':
        return `Nivel >= L${admittedCategories.min_level}`
      case 'LESS_THAN':
        return `Nivel <= L${admittedCategories.max_level}`
      case 'BETWEEN':
        return `Nivel L${admittedCategories.min_level} - L${admittedCategories.max_level}`
      default:
        return 'Cualquier nivel'
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'PENDING':
        return 'warning'
      case 'CONFIRMED':
        return 'success'
      case 'CANCELLED':
        return 'error'
      case 'EXPIRED':
        return 'default'
      default:
        return 'default'
    }
  }

  const getStatusText = (status: string) => {
    switch (status) {
      case 'PENDING':
        return 'Pendiente'
      case 'CONFIRMED':
        return 'Confirmado'
      case 'CANCELLED':
        return 'Cancelado'
      case 'EXPIRED':
        return 'Expirado'
      default:
        return status
    }
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
            <EventIcon sx={{ fontSize: 64 }} />
            <Typography variant="h3" component="h1" align="center" fontWeight={700}>
              Partidos Disponibles
            </Typography>
            <Typography variant="h6" align="center" sx={{ opacity: 0.95, maxWidth: 600 }}>
              Encuentra equipos que buscan rivales para jugar desde hoy en adelante
            </Typography>
          </Stack>
        </Paper>

        {/* Loading State */}
        {loading && (
          <Box display="flex" justifyContent="center" py={4}>
            <CircularProgress />
          </Box>
        )}

        {/* Error State */}
        {error && (
          <Alert severity="error" onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        {/* Empty State */}
        {!loading && !error && announcements.length === 0 && (
          <Paper sx={{ p: 4, textAlign: 'center' }}>
            <EventIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" color="text.secondary" gutterBottom>
              No hay partidos disponibles
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Sé el primero en publicar un partido
            </Typography>
          </Paper>
        )}

        {/* Announcements List */}
        {!loading && !error && announcements.length > 0 && (
          <Grid container spacing={3}>
            {announcements.map((announcement) => (
              <Grid item xs={12} md={6} key={announcement.id}>
                <Card elevation={2} sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                  <CardContent sx={{ flexGrow: 1 }}>
                    <Stack spacing={2}>
                      {/* Header */}
                      <Box display="flex" justifyContent="space-between" alignItems="start">
                        <Box>
                          <Typography variant="h6" fontWeight={700} gutterBottom>
                            {announcement.team_name}
                          </Typography>
                          <Chip
                            label={announcement.sport}
                            icon={<SportsIcon />}
                            size="small"
                            color="primary"
                            variant="outlined"
                          />
                        </Box>
                        <Chip
                          label={getStatusText(announcement.status)}
                          size="small"
                          color={getStatusColor(announcement.status) as any}
                        />
                      </Box>

                      <Divider />

                      {/* Date and Time */}
                      <Box>
                        <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 1 }}>
                          <EventIcon fontSize="small" color="action" />
                          <Typography variant="body2" color="text.secondary" fontWeight={600}>
                            {formatDate(announcement.day)}
                          </Typography>
                        </Stack>
                        <Stack direction="row" spacing={1} alignItems="center">
                          <AccessTimeIcon fontSize="small" color="action" />
                          <Typography variant="body2" color="text.secondary">
                            {formatTime(announcement.time_slot.start_time)} - {formatTime(announcement.time_slot.end_time)}
                          </Typography>
                        </Stack>
                      </Box>

                      {/* Location */}
                      <Box>
                        <Stack direction="row" spacing={1} alignItems="center">
                          <LocationOnIcon fontSize="small" color="action" />
                          <Typography variant="body2" color="text.secondary">
                            {announcement.location.locality}, {announcement.location.province}, {announcement.location.country}
                          </Typography>
                        </Stack>
                      </Box>

                      {/* Categories */}
                      <Box>
                        <Stack direction="row" spacing={1} alignItems="center">
                          <GroupsIcon fontSize="small" color="action" />
                          <Typography variant="body2" color="text.secondary">
                            {getCategoryText(announcement.admitted_categories)}
                          </Typography>
                        </Stack>
                      </Box>
                    </Stack>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        )}
      </Stack>
    </Box>
  )
}

