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
import { CURRENT_ACCOUNT_ID } from '../../../../shared/constants/session'

export function MyAnnouncementsPage() {
  const { findAccountMatchAnnouncementsUseCase } = useMatchAnnouncementContext()
  const [announcements, setAnnouncements] = useState<MatchAnnouncement[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    findAccountMatchAnnouncementsUseCase.execute(CURRENT_ACCOUNT_ID).then((result) => {
      if (result.success) {
        setAnnouncements(result.announcements)
      } else {
        setError(result.error ?? 'Error al cargar los anuncios')
      }
      setLoading(false)
    })
  }, [])

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
        return `Categorías: ${admittedCategories.categories?.map((c) => `L${c}`).join(', ')}`
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
      case 'PENDING': return 'warning'
      case 'ACCEPTED': return 'success'
      case 'CANCELLED': return 'error'
      default: return 'default'
    }
  }

  const getStatusText = (status: string) => {
    switch (status) {
      case 'PENDING': return 'Pendiente'
      case 'ACCEPTED': return 'Aceptado'
      case 'CANCELLED': return 'Cancelado'
      default: return status
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
            <EventIcon sx={{ fontSize: 64 }} />
            <Typography variant="h3" component="h1" align="center" fontWeight={700}>
              Mis Anuncios
            </Typography>
            <Typography variant="h6" align="center" sx={{ opacity: 0.95, maxWidth: 600 }}>
              Partidos que publicaste
            </Typography>
          </Stack>
        </Paper>

        {loading && (
          <Box display="flex" justifyContent="center" py={4}>
            <CircularProgress />
          </Box>
        )}

        {!loading && error && (
          <Alert severity="error">{error}</Alert>
        )}

        {!loading && !error && announcements.length === 0 && (
          <Paper sx={{ p: 4, textAlign: 'center' }}>
            <EventIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" color="text.secondary">
              No publicaste ningún partido todavía
            </Typography>
          </Paper>
        )}

        {!loading && !error && announcements.length > 0 && (
          <Grid container spacing={3}>
            {announcements.map((announcement) => (
              <Grid item xs={12} sm={6} md={4} key={announcement.id}>
                <Card elevation={2} sx={{ height: '100%' }}>
                  <CardContent>
                    <Stack spacing={2}>
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

                      <Box>
                        <Stack direction="row" spacing={1} alignItems="center">
                          <LocationOnIcon fontSize="small" color="action" />
                          <Typography variant="body2" color="text.secondary">
                            {announcement.location.locality}, {announcement.location.province}
                          </Typography>
                        </Stack>
                      </Box>

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
