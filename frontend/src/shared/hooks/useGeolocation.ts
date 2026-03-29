import { useState, useCallback } from 'react'

export type GeoStatus = 'idle' | 'loading' | 'granted' | 'denied' | 'unavailable'

export interface GeolocationState {
  status: GeoStatus
  latitude: number | null
  longitude: number | null
}

export interface UseGeolocationReturn extends GeolocationState {
  requestLocation: () => void
  reset: () => void
}

export function useGeolocation(): UseGeolocationReturn {
  const [state, setState] = useState<GeolocationState>({
    status: 'idle',
    latitude: null,
    longitude: null,
  })

  const requestLocation = useCallback(() => {
    if (!('geolocation' in navigator)) {
      setState({ status: 'unavailable', latitude: null, longitude: null })
      return
    }

    setState((prev) => ({ ...prev, status: 'loading' }))

    navigator.geolocation.getCurrentPosition(
      (position) => {
        setState({
          status: 'granted',
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
        })
      },
      (error) => {
        const status = error.code === GeolocationPositionError.PERMISSION_DENIED ? 'denied' : 'unavailable'
        setState({ status, latitude: null, longitude: null })
      },
      { timeout: 10_000, enableHighAccuracy: false }
    )
  }, [])

  const reset = useCallback(() => {
    setState({ status: 'idle', latitude: null, longitude: null })
  }, [])

  return { ...state, requestLocation, reset }
}
