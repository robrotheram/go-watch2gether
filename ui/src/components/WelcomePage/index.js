import './home.less';

import logo from './logo.jpg'

import { useEffect, useState } from "react"
import { Layout, Button, Divider, Typography, Checkbox } from 'antd';
import { Form, Input } from 'antd';
import { Alert } from 'antd';
import { connect } from 'react-redux'

import { withRouter } from "react-router";
import queryString from 'query-string'
import {PageFooter} from '../common/PageFooter'

import store , {BASE_URL} from '../../store'
import {checklogin} from '../../store/user/user.actions'
import {ROOM_ERROR} from '../../store/room/room.types'
import { join, leave, clearError, getMeta } from '../../store/room/room.actions'

const { Title, Paragraph, Text, Link } = Typography;
const { Content } = Layout;


function Home(props) {

  const [form] = Form.useForm();
  const {name, id} = props.room
  const [discord_login, setLoginURL] = useState(BASE_URL+"/auth/login")

  const layout = {
    labelCol: { span: 6 },
    wrapperCol: { span: 14 },
  };
  const tailLayout = {
    wrapperCol: { offset: 6, span: 14 },
  };

  useEffect(() => {
    const values = queryString.parse(props.location.search);
    if (values.room !== undefined) {
      setLoginURL(discord_login+"?next=/?room="+values.room)
    }
    
    props.checklogin(values.room);
  }, []);


  const onFinish = values => {
    console.log("Valuse", values, id)
    if (values.roomname === name){
      props.join(id, values.roomname, values.username, values.anonymous)
    } else {
      props.join("", values.roomname, values.username, values.anonymous)
    }
  };

  const onFinishFailed = errorInfo => {
    console.log('Failed:', errorInfo);
  };

  const handleClose = () => {
    props.clearError();
  };



  useEffect(() => {
    const values = queryString.parse(props.location.search);
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
      props.getMeta(values.room)
    }
  }, [props.location.search]);
  
  
  useEffect(() => {
    form.setFieldsValue({ roomname: props.room.name });
    form.setFieldsValue({ username: props.auth.username });
  }, [props.room, props.auth]);
 

  return (
    <div className="wrap-login">
        <Content className="login-form">
            <Typography>
            <div style={{"width": "500px", marginBottom:"70px"}}>
              <img src={logo} alt="watch2gether logo" style={{"float": "left", "width": "80px", marginRight:"30px"}} />
              <Title style={{width: "400px",  display:"block", paddingTop:"16px"}} level={1}>Watch2Gether</Title>
            </div> 
            </Typography>

            {props.error !== "" ? (
              <Alert
                message="Error"
                description={props.error}
                type="error"
                showIcon
                closable
                afterClose={handleClose}
                style={{ "marginBottom": "20px" }}
              />
            ) : null}
            
           {/* <Form
              {...layout}
              name="basic"
              form={form}
              initialValues={{ anonymous: true }}
              onFinish={onFinish}
              onFinishFailed={onFinishFailed}
            >
                <Form.Item
                label="Room id"
                name="roomid"
               >
                <Input size="large" />
              </Form.Item> 

              <Form.Item
                label="Room Name"
                name="roomname"
                rules={[
                  { required: true, message: 'Please input your room Name!' }, 
                  { min: 4, message: 'Room name must be minimum 4 characters.'},
                  {
                    pattern: new RegExp(
                      /^[a-zA-Z0-9@~`!@#$%^&*()_=+\\\\';:"\\/?>.<,-]+$/i
                    ),
                    message: "Valid characters are letters, numbers"
                  }
                
                ]}
              >
                <Input size="large" />
              </Form.Item>
              <Form.Item
                label="Username"
                name="username"
                style={{ "marginTop": "20px" }}
                
              >
                <Input size="large" disabled={props.auth.auth}/>
              </Form.Item>

              {!props.auth.auth ? 
              <div>
              <Form.Item {...tailLayout} name="anonymous" valuePropName="checked">
                <Checkbox disabled>Be anonymous </Checkbox>
              </Form.Item>

              <Form.Item {...tailLayout}>
                <Button size="large" type="primary" htmlType="submit" style={{ "width": "100%", marginTop: "20px" }}>
                  Login
                </Button>

                <Divider/>

                <Button href={BASE_URL+"/auth/login"} size="large" type="primary" style={{ "width": "100%", marginTop: "0px", backgroundColor: "#7289da", border: "none" }}>
                  Login with Discord
                </Button>
              </Form.Item>
              </div>
              :<div>
                
              <Form.Item {...tailLayout}>
                <Button size="large" type="primary" htmlType="submit" style={{ "width": "100%", marginTop: "20px" }}>
                  Join Room
                </Button>
              </Form.Item> 
              </div> 
              }


            </Form> */}

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
            <Button href={discord_login} size="large" type="primary" style={{ "width": "100%", marginTop: "0px", backgroundColor: "#7289da", border: "none" }}>
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


