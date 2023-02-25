import React, { useState } from 'react'
import ReactDOM from 'react-dom/client'
import './index.css'
import { ConfigProvider, theme } from 'antd';
import RoomPage from './components/RoomsPage/RoomPage';
import Home from './components/WelcomePage/Home';
import { BrowserRouter, Navigate, Route, Routes, useLocation, useNavigate } from 'react-router-dom';
import { UserContextProvider } from './context/UserContext';
import { PlayerApp } from './components/PlayerApp/PlayerApp';
import { Empty } from 'antd';
import { RoomContextProvider } from './context/RoomContext';

const PrivateRoute = ({ authed, children }) => {
  let location = useLocation();
  const { pathname } = location;
  const lastItem = pathname.substring(pathname.lastIndexOf('/') + 1);
  console.log("REST", authed)
  // if (!authed) {
  //   return <Navigate to={`/?room=${lastItem}`} replace />;
  // }
  return (
    children
  );
};

const Router = () => {
  return (
    <BrowserRouter>
      <UserContextProvider>
        <Routes>
          <Route path="/app" element={<RoomPage />}>
            <Route path="/app/room/:id" element={<PlayerApp />}/>
            <Route path='/app' element={
              <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={(<span>Sorry Could not find a romm please select one</span>)} style={{ paddingTop: '250px' }} />
            } />
          </Route>
          <Route path="/" element={<Home />} />
        </Routes>
      </UserContextProvider>
    </BrowserRouter>
  )
}
ReactDOM.createRoot(document.getElementById('root')).render(
  <ConfigProvider
    theme={{
      algorithm: [theme.darkAlgorithm],
      components: {
        Layout: {
          colorBgHeader: "rgb(20, 20, 20)",
          colorBgTrigger: "rgb(20, 20, 20)"
        },
        Menu: {
          colorItemBg: "rgb(20, 20, 20)"
        }
      }
    }}>
    <Router />
  </ConfigProvider>,
)
