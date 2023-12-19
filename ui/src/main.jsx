// import React from 'react'
// import ReactDOM from 'react-dom/client'
// import { App } from './app/app'
// import './main.css'

// ReactDOM.createRoot(document.getElementById('root')).render(
//   <React.StrictMode>
//     <App/>
//   </React.StrictMode>,
// )

import React from 'react'
import ReactDOM from 'react-dom/client'
import Index from './pages/index'
import App from './pages/app'
import './main.css'
import { BrowserRouter, Route, Routes} from 'react-router-dom';
import { AppController } from './pages/app/controller';
import { ErrroPage } from './pages/app/error'
import PlaylistPage from './pages/app/playlist'



const Router = () => {
  return (
    <BrowserRouter>
        <Routes>
          <Route path="app" element={<App />}>
            <Route path=":id" element={<AppController />}/>
            <Route path=":id/playlists" element={<PlaylistPage />}/>
            <Route index element={<ErrroPage/>}/>
          </Route>
          <Route path="/" element={<Index />} />
        </Routes>
    </BrowserRouter>
  )
}

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <Router />
  </React.StrictMode>,
)