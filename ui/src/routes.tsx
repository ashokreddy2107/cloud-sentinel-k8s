import { createBrowserRouter, Navigate, useLocation } from 'react-router-dom'

import App from './App'
import { InitCheckRoute } from './components/init-check-route'
import { ProtectedRoute } from './components/protected-route'
import { getSubPath } from './lib/subpath'
import { CRListPage } from './pages/cr-list-page'
import { InitializationPage } from './pages/initialization'
import { LoginPage } from './pages/login'
import { Overview } from './pages/overview'
import { ResourceDetail } from './pages/resource-detail'
import { ResourceList } from './pages/resource-list'
import { SettingsPage } from './pages/settings'
import { useCluster } from './hooks/use-cluster'

const subPath = getSubPath()

function RootRedirector() {
  const { currentCluster } = useCluster()
  // Wait for cluster to be resolved
  if (!currentCluster) {
    return null // generic loading or let the App's loading state handle it
  }
  return <Navigate to={`/c/${currentCluster}/dashboard`} replace />
}

function ClusterRedirector() {
  const { currentCluster } = useCluster()
  const location = useLocation()

  if (!currentCluster) {
    return null
  }

  // Preserve the current path but prepend cluster
  // e.g. /pods -> /c/dev/pods
  return <Navigate to={`/c/${currentCluster}${location.pathname}`} replace />
}

export const router = createBrowserRouter(
  [
    {
      path: '/setup',
      element: <InitializationPage />,
    },
    {
      path: '/login',
      element: (
        <InitCheckRoute>
          <LoginPage />
        </InitCheckRoute>
      ),
    },
    {
      path: '/',
      element: (
        <InitCheckRoute>
          <ProtectedRoute>
            <App />
          </ProtectedRoute>
        </InitCheckRoute>
      ),
      children: [
        {
          index: true,
          element: <RootRedirector />,
        },
        {
          path: 'settings',
          element: <SettingsPage />,
        },
        {
          path: 'c/:cluster',
          children: [
            {
              index: true,
              element: <Navigate to="dashboard" replace />,
            },
            {
              path: 'dashboard',
              element: <Overview />,
            },
            {
              path: 'crds/:crd',
              element: <CRListPage />,
            },
            {
              path: 'crds/:resource/:namespace/:name',
              element: <ResourceDetail />,
            },
            {
              path: 'crds/:resource/:name',
              element: <ResourceDetail />,
            },
            {
              path: ':resource/:name',
              element: <ResourceDetail />,
            },
            {
              path: ':resource',
              element: <ResourceList />,
            },
            {
              path: ':resource/:namespace/:name',
              element: <ResourceDetail />,
            },
          ],
        },
        {
          // Catch-all for legacy/absolute paths that forgot the cluster prefix
          path: '*',
          element: <ClusterRedirector />
        }
      ],
    },
  ],
  {
    basename: subPath,
  }
)
