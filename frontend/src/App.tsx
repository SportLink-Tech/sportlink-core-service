import { Routes, Route, Navigate } from 'react-router-dom'
import { Layout } from './components/Layout'
import { MatchOfferModule } from './features/matchoffer/MatchOfferModule'
import { ListMatchOffersPage } from './features/matchoffer/ui/pages/ListMatchOffersPage'
import { CreateMatchOfferPage } from './features/matchoffer/ui/pages/CreateMatchOfferPage'
import { TeamModule } from './features/team/TeamModule'
import { CreateTeamPage } from './features/team/ui/pages/CreateTeamPage'
import { MyTeamsPage } from './features/team/ui/pages/MyTeamsPage'
import { MyOffersPage } from './features/matchoffer/ui/pages/MyOffersPage'
import { MatchRequestModule } from './features/matchrequest/MatchRequestModule'
import { LoginPage } from './features/auth/ui/pages/LoginPage'
import { ProfilePage } from './features/auth/ui/pages/ProfilePage'
import { useAuth } from './features/auth/context/AuthContext'

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const { accountId } = useAuth()
  if (!accountId) return <Navigate to="/login" replace />
  return <>{children}</>
}

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        path="/*"
        element={
          <PrivateRoute>
            <MatchOfferModule>
              <MatchRequestModule>
                <TeamModule>
                  <Layout>
                    <Routes>
                      <Route path="/" element={<ListMatchOffersPage />} />
                      <Route path="/create" element={<CreateMatchOfferPage />} />
                      <Route path="/create-team" element={<CreateTeamPage />} />
                      <Route path="/my-teams" element={<MyTeamsPage />} />
                      <Route path="/my-offers" element={<MyOffersPage />} />
                      <Route path="/profile" element={<ProfilePage />} />
                    </Routes>
                  </Layout>
                </TeamModule>
              </MatchRequestModule>
            </MatchOfferModule>
          </PrivateRoute>
        }
      />
    </Routes>
  )
}

export default App
