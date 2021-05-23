import './home.less';

import logo from './logo.jpg'

import { useEffect, useState } from "react"
import { Layout, Button, Typography } from 'antd';
import { Alert } from 'antd';
import { connect } from 'react-redux'

import { withRouter } from "react-router";
import queryString from 'query-string'
import {PageFooter} from '../common/PageFooter'

import store , {BASE_URL} from '../../store'
import {checklogin} from '../../store/user/user.actions'
import {ROOM_ERROR} from '../../store/room/room.types'
import { join, leave, clearError, getMeta } from '../../store/room/room.actions'
import axios from 'axios';

const { Title, Paragraph } = Typography;
const { Content } = Layout;


const Home = ({location, checklogin, clearError, getMeta, error}) => {
  const [discord_login, setLoginURL] = useState(BASE_URL+"/auth/login")
  const [botid, setBot] = useState("")

  useEffect(() => {
    const values = queryString.parse(location.search);
    if (values.room !== undefined) {
      setLoginURL(discord_login+"?next=/?room="+values.room)
    }
    
    checklogin(values.room);
  }, [location,discord_login, checklogin ]);

  const handleClose = () => {
    clearError();
  };

  useEffect(() => {
    axios.get(BASE_URL+"/config").then(res => {
      setBot(res.data.bot)
    })
  }, []);

  const inviteBotUrl = (bot) => {
    return "https://discord.com/oauth2/authorize?client_id="+bot+"&scope=bot"
  }

  useEffect(() => {
    const values = queryString.parse(location.search);
    const err = values.error
    console.log("QUETR", values.room);

    console.log("Error", err);
    
    if (err !== undefined && err !== "") {
      store.dispatch( {
        type: ROOM_ERROR,
        error:err,
      })
    }
    if(values.room !== undefined && values.room !== ""){
      getMeta(values.room)
    }
  }, [getMeta, location.search]);

  return (
    <div className="wrap-login">
        {botid !== "" ?
        <Button target="_blank" href={inviteBotUrl(botid)} size="large" type="primary" shape="round" style={{ 
          position:"fixed",
          top:"20px",
          right:"20px",
          marginTop: "0px",
          padding: "0px 20px",
          backgroundColor: "#7289da",
          border: "none" }}
        > 
          Add the Discord Bot 
        </Button>
        : null }
        
        <Content className="login-form">
            <Typography>
            <div style={{"width": "500px", marginBottom:"70px"}}>
              <img src={logo} alt="watch2gether logo" style={{"float": "left", "width": "80px", marginRight:"30px"}} />
              <Title style={{width: "400px",  display:"block", paddingTop:"16px"}} level={1}>Watch2Gether</Title>
            </div> 
            </Typography>

            {error !== "" ? (
              <Alert
                message="Error"
                description={error}
                type="error"
                showIcon
                closable
                afterClose={handleClose}
                style={{ "marginBottom": "20px" }}
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
            <Button href={discord_login} size="large" shape="round" type="primary" style={{ padding: "0px 20px", "width": "100%", marginTop: "0px", backgroundColor: "#7289da", border: "none" }}>
                  Login with Discord
            </Button>

            <PageFooter style={{ textAlign: 'center', position: "absolute", bottom:"0px", left:"0px", width:"560px", height:"50px", padding: "15px 50px 28px 50px" }}/>
        </Content>
      </div>
    
  );
}

const mapStateToProps = (state) => {
  console.log(state)
  return ({ error: state.room.error, room: state.room, auth: state.user })
}
export default withRouter(connect(mapStateToProps, {checklogin, join, leave, clearError, getMeta })(Home))


