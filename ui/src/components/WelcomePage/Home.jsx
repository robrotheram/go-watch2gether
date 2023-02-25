import './home.less';
import React, { useContext } from "react"
import { useState } from 'react';
import {Layout, Button, Typography, Alert, Spin} from 'antd';

import logo from './logo.jpg';
import { UserContext } from '../../context/UserContext';
import { BASE_URL } from '../../context/config';

const { Title } = Typography;
const { Content } = Layout;

const Home = () => {
  const [discord_login, setLoginURL] = useState(`${BASE_URL}/auth/login`);
  const [user, loading] = useContext(UserContext);
  const [error, setError] = useState("");

  return (
    <div className="wrap-login">
      <Content className="login-form">
        <Typography>
          <div className="welcomeHeading">
            <img src={logo} alt="watch2gether logo" className="logo" />
            <Title className="title" level={1} style={{ margin: "10px" }}>Watch2Gether</Title>
          </div>
        </Typography>
       
        {error !== '' ? (
          <Alert
            message="There was a problem logging you in"
            description={error}
            type="error"
            showIcon
            closable
          />
        ) : null}
        <br />
        {loading?
          <Spin size="large"  style={{width: '100%', marginTop: '0px', border: 'none', color:"#FFF"}}/>:
          <Button
            href={discord_login}
            size="large"
            shape="round"
            type="primary"
            style={{
              width: '100%', marginTop: '0px', backgroundColor: '#7289da', border: 'none',
            }}
          >
            Login with Discord
          </Button>
        }

      </Content>
    </div>

  );
};

export default Home
// const mapStateToProps = (state) => {
//   console.log(state);
//   return ({ error: state.room.error, room: state.room, auth: state.user });
// };
// export default withRouter(connect(mapStateToProps, {
//   checklogin, join, leave, clearError, getMeta,
// })(Home));
