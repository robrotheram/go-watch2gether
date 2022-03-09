import './home.less';
import React from "react"
import { useEffect, useState } from 'react';
import {
  Layout, Button, Typography, Alert,
} from 'antd';

import { connect } from 'react-redux';

import { withRouter } from 'react-router';
import queryString from 'query-string';
import axios from 'axios';
import { PageFooter } from '../common/PageFooter';

import store, { BASE_URL } from '../../store';
import { checklogin } from '../../store/user/user.actions';
import { ROOM_ERROR } from '../../store/event.types';
import {
  join, leave, clearError, getMeta,
} from '../../store/room/room.actions';
import logo from './logo.jpg';

const { Title, Paragraph } = Typography;
const { Content } = Layout;

const Home = ({
  location, checklogin, clearError, getMeta, error,
}) => {
  const [discord_login, setLoginURL] = useState(`${BASE_URL}/auth/login`);
  const [botid, setBot] = useState('');

  useEffect(() => {
    const values = queryString.parse(location.search);
    if (values.room !== undefined) {
      setLoginURL(`${BASE_URL}/auth/login?next=/?room=${values.room}`);
    }
    checklogin(values.room);
  }, [location.search, checklogin]);

  const handleClose = () => {
    clearError();
  };

  useEffect(() => {
    axios.get(`${BASE_URL}/config`).then((res) => {
      setBot(res.data.bot);
    });
  }, []);

  const inviteBotUrl = (bot) => `https://discord.com/oauth2/authorize?client_id=${bot}&permissions=0&scope=bot%20applications.commands`;
  useEffect(() => {
    const values = queryString.parse(location.search);
    const err = values.error;
    console.log('QUETR', values.room);

    console.log('Error', err);

    if (err !== undefined && err !== '') {
      store.dispatch({
        type: ROOM_ERROR,
        error: err,
      });
    }
    if (values.room !== undefined && values.room !== '') {
      getMeta(values.room);
    }
  }, [getMeta, location.search]);

  return (
    <div className="wrap-login">
      {botid !== '' ? (
        <Button target="_blank" href={inviteBotUrl(botid)} size="large" type="primary" shape="round" className="discordBotButton">
          Add the Discord Bot
        </Button>
      )
        : null }

      <Content className="login-form">
        <Typography>
          <div className="welcomeHeading">
            <img src={logo} alt="watch2gether logo" className="logo" />
            <Title className="title" level={1}>Watch2Gether</Title>
          </div>
        </Typography>

        {error !== '' ? (
          <Alert
            message="Error"
            description={error}
            type="error"
            showIcon
            closable
            afterClose={handleClose}
            style={{ marginBottom: '20px' }}
          />
        ) : null}

        <Typography>
          <Paragraph>
            Ever wanted to watch youtube videos in-sync with your friends, via. web-browser? or mp4s?
          </Paragraph>
          <Paragraph>
            Its yet another video sync website it currently support Youtube,and Videos hosted on your own fileserver that you totally legally own 😉
          </Paragraph>
          <Paragraph>
            Also comes with a Discord Bot, Playlist support and fun!
          </Paragraph>
        </Typography>
        <Button
          href={discord_login}
          size="large"
          shape="round"
          type="primary"
          style={{
            padding: '0px 20px', width: '100%', marginTop: '0px', backgroundColor: '#7289da', border: 'none',
          }}
        >
          Login with Discord
        </Button>

        <PageFooter className="welcomeFooter" />
      </Content>
    </div>

  );
};

const mapStateToProps = (state) => {
  console.log(state);
  return ({ error: state.room.error, room: state.room, auth: state.user });
};
export default withRouter(connect(mapStateToProps, {
  checklogin, join, leave, clearError, getMeta,
})(Home));
