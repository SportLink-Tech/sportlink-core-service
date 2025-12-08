import { useState, useEffect, useMemo, useCallback } from 'react'
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
  Button,
} from '@mui/material'
import EventIcon from '@mui/icons-material/Event'
import LocationOnIcon from '@mui/icons-material/LocationOn'
import SportsIcon from '@mui/icons-material/Sports'
import GroupsIcon from '@mui/icons-material/Groups'
import AccessTimeIcon from '@mui/icons-material/AccessTime'
import { useMatchAnnouncementContext } from '../../context/MatchAnnouncementContext'
import { MatchAnnouncement } from '../../../../shared/types/matchAnnouncement'
import { MatchAnnouncementFilters } from '../components/MatchAnnouncementFilters'

const ITEMS_PER_PAGE = 9 // 3 filas x 3 columnas

export function ListMatchAnnouncementsPage() {
  const { findMatchAnnouncementsUseCase } = useMatchAnnouncementContext()

  const [announcements, setAnnouncements] = useState<MatchAnnouncement[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(0)

  // Filter state
  const [selectedSports, setSelectedSports] = useState<string[]>([])
  const [fromDate, setFromDate] = useState<string>('')
  const [toDate, setToDate] = useState<string>('')

  const loadAnnouncements = useCallback(async () => {
    setLoading(true)
    setError(null)

    const offset = (currentPage - 1) * ITEMS_PER_PAGE
    const limit = ITEMS_PER_PAGE

    const query: any = {
      limit,
      offset,
    }

    // Add filters if they are set
    if (selectedSports.length > 0) {
      query.sports = selectedSports
    }
    if (fromDate) {
      query.fromDate = fromDate
    }
    if (toDate) {
      query.toDate = toDate
    }

    const result = await findMatchAnnouncementsUseCase.execute(query)

    if (result.success) {
      setAnnouncements(result.announcements)
      setTotalPages(result.pagination.outOf)
      // Asegurarnos de que currentPage esté sincronizado con el backend
      if (result.pagination.number !== currentPage) {
        setCurrentPage(result.pagination.number)
      }
    } else {
      setError(result.error || 'Error al cargar los anuncios')
      setAnnouncements([])
      setTotalPages(0)
    }

    setLoading(false)
  }, [currentPage, findMatchAnnouncementsUseCase, selectedSports, fromDate, toDate])

  useEffect(() => {
    loadAnnouncements()
  }, [loadAnnouncements])

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

  // Calcular qué números de página mostrar basándose en el total de páginas del backend
  const getVisiblePageNumbers = () => {
    const maxVisible = 10
    
    if (totalPages === 0) {
      return []
    }

    const pages: number[] = []

    // Calcular rango de páginas a mostrar
    let startPage = Math.max(1, currentPage - 4)
    let endPage = Math.min(totalPages, currentPage + 5)

    // Ajustar si estamos cerca del inicio
    if (currentPage <= 5) {
      startPage = 1
      endPage = Math.min(maxVisible, totalPages)
    }

    // Ajustar si estamos cerca del final
    if (currentPage > totalPages - 5) {
      startPage = Math.max(1, totalPages - maxVisible + 1)
      endPage = totalPages
    }

    // Generar números de página
    for (let i = startPage; i <= endPage; i++) {
      pages.push(i)
    }

    return pages
  }

  const visiblePageNumbers = useMemo(() => getVisiblePageNumbers(), [currentPage, totalPages])

  const handlePageChange = (page: number) => {
    // Solo permitir cambiar a páginas válidas (no duplicadas y dentro del rango válido)
    if (page >= 1 && page !== currentPage) {
      setCurrentPage(page)
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }
  }

  const handleNextPage = () => {
    if (currentPage < totalPages) {
      setCurrentPage(currentPage + 1)
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }
  }

  const handleFiltersChange = (filters: { sports: string[]; fromDate: string; toDate: string }) => {
    setSelectedSports(filters.sports)
    setFromDate(filters.fromDate)
    setToDate(filters.toDate)
    // Reset to first page when filters change
    setCurrentPage(1)
  }

  const shouldShowPagination = totalPages > 1

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

        {/* Filters Bar - Debajo del banner */}
        <MatchAnnouncementFilters
          selectedSports={selectedSports}
          fromDate={fromDate}
          toDate={toDate}
          onFiltersChange={handleFiltersChange}
        />

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
          <Stack spacing={3}>
            <Grid container spacing={3}>
              {announcements.map((announcement) => (
                <Grid item xs={12} sm={6} md={4} lg={4} key={announcement.id}>
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

            {/* Pagination */}
            {shouldShowPagination && (
              <Box display="flex" justifyContent="center" alignItems="center" gap={1} flexWrap="wrap" sx={{ mt: 2 }}>
                {visiblePageNumbers.map((pageNum, index) => (
                  <Button
                    key={index}
                    variant={pageNum === currentPage ? 'contained' : 'outlined'}
                    onClick={() => typeof pageNum === 'number' && handlePageChange(pageNum)}
                    disabled={typeof pageNum !== 'number'}
                    sx={{
                      minWidth: 40,
                      height: 40,
                      borderRadius: 1,
                      fontWeight: pageNum === currentPage ? 700 : 400,
                      borderColor: pageNum === currentPage ? 'primary.main' : 'divider',
                      '&:hover': {
                        borderColor: 'primary.main',
                        backgroundColor: pageNum === currentPage ? 'primary.dark' : 'action.hover',
                      },
                    }}
                  >
                    {pageNum}
                  </Button>
                ))}
                {currentPage < totalPages && (
                  <Button
                    variant="outlined"
                    onClick={handleNextPage}
                    disabled={currentPage >= totalPages}
                    sx={{
                      minWidth: 100,
                      height: 40,
                      borderRadius: 1,
                      ml: 1,
                      textTransform: 'none',
                    }}
                  >
                    Siguiente &gt;
                  </Button>
                )}
              </Box>
            )}
          </Stack>
        )}
      </Stack>
    </Box>
  )
}

