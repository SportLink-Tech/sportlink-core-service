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
import { MySentRequestsPage } from './features/matchrequest/ui/pages/MySentRequestsPage'
import { MyReceivedRequestsPage } from './features/matchrequest/ui/pages/MyReceivedRequestsPage'
import { MyMatchesPage } from './features/match/ui/pages/MyMatchesPage'
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
                      <Route path="/my-offers/new" element={<CreateMatchOfferPage />} />
                      <Route path="/teams" element={<MyTeamsPage />} />
                      <Route path="/teams/new" element={<CreateTeamPage />} />
                      <Route path="/my-offers" element={<MyOffersPage />} />
                      <Route path="/my-requests/sent" element={<MySentRequestsPage />} />
                      <Route path="/my-requests/received" element={<MyReceivedRequestsPage />} />
                      <Route path="/my-matches" element={<MyMatchesPage />} />
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
