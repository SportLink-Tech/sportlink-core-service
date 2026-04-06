import { Routes, Route } from 'react-router-dom'
import { Layout } from './components/Layout'
import { MatchOfferModule } from './features/matchoffer/MatchOfferModule'
import { ListMatchOffersPage } from './features/matchoffer/ui/pages/ListMatchOffersPage'
import { CreateMatchOfferPage } from './features/matchoffer/ui/pages/CreateMatchOfferPage'
import { TeamModule } from './features/team/TeamModule'
import { CreateTeamPage } from './features/team/ui/pages/CreateTeamPage'
import { MyTeamsPage } from './features/team/ui/pages/MyTeamsPage'
import { MyOffersPage } from './features/matchoffer/ui/pages/MyOffersPage'
import { MatchRequestModule } from './features/matchrequest/MatchRequestModule'

/**
 * App Component
 * Wraps the application with MatchOfferModule for Dependency Injection
 * Following Atomic Hexagonal Architecture
 */
function App() {
  return (
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
            </Routes>
          </Layout>
        </TeamModule>
      </MatchRequestModule>
    </MatchOfferModule>
  )
}

export default App

