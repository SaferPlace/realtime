import React from 'react'
import ReactDOM from 'react-dom/client'
import {
  createBrowserRouter,
  RouterProvider
} from 'react-router-dom'
import Root from './routes/root'
import { CssBaseline } from '@mui/material'
import Login from './routes/login'
import './i18n'
import Home from './routes/home'
import IncidentList from './routes/incident/list'
import Incident from './routes/incident/single'
import { incidentLoader, incidentsInRadiusLoader } from './routes/incident/loaders'

const router = createBrowserRouter([
  {
    path: '/',
    element: <Root />,
    children: [
      {
        path: '/',
        element: <Home />,
      }, {
        path: 'incidents',
        loader: incidentsInRadiusLoader,
        Component: IncidentList,
      }, {
        path: 'incident/:id',
        loader: incidentLoader,
        Component: Incident,
      }
    ]
  }, {
    path: '/login',
    element: <Login />,
  }
])

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
    <CssBaseline />
  </React.StrictMode>,
)
