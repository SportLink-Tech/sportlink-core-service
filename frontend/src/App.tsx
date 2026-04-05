import { Routes, Route } from 'react-router-dom'
import { Layout } from './components/Layout'
import { MatchAnnouncementModule } from './features/matchannouncement/MatchAnnouncementModule'
import { ListMatchAnnouncementsPage } from './features/matchannouncement/ui/pages/ListMatchAnnouncementsPage'
import { CreateMatchAnnouncementPage } from './features/matchannouncement/ui/pages/CreateMatchAnnouncementPage'
import { TeamModule } from './features/team/TeamModule'
import { CreateTeamPage } from './features/team/ui/pages/CreateTeamPage'
import { MyTeamsPage } from './features/team/ui/pages/MyTeamsPage'
import { MatchRequestModule } from './features/matchrequest/MatchRequestModule'

/**
 * App Component
 * Wraps the application with MatchAnnouncementModule for Dependency Injection
 * Following Atomic Hexagonal Architecture
 */
function App() {
  return (
    <MatchAnnouncementModule>
      <MatchRequestModule>
        <TeamModule>
          <Layout>
            <Routes>
              <Route path="/" element={<ListMatchAnnouncementsPage />} />
              <Route path="/create" element={<CreateMatchAnnouncementPage />} />
              <Route path="/create-team" element={<CreateTeamPage />} />
              <Route path="/my-teams" element={<MyTeamsPage />} />
            </Routes>
          </Layout>
        </TeamModule>
      </MatchRequestModule>
    </MatchAnnouncementModule>
  )
}

export default App

