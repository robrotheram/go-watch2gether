import './home.less';

import logo from './logo.jpg'

import { useEffect } from "react"
import { Layout, Button, Divider, Image, Typography, Checkbox } from 'antd';
import { Form, Input } from 'antd';
import { Alert } from 'antd';
import { connect } from 'react-redux'
import { join, leave, clearError, getMeta } from '../store/room/room.actions'
import { withRouter } from "react-router";
import queryString from 'query-string'
import {PageFooter} from '../components/PageFooter'


const { Content } = Layout;
const { Title, Paragraph, Text } = Typography;


function Home(props) {

  const [form] = Form.useForm();

  const {name, id} = props.room

  const layout = {
    labelCol: { span: 6 },
    wrapperCol: { span: 14 },
  };
  const tailLayout = {
    wrapperCol: { offset: 6, span: 14 },
  };


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
    console.log("QUETR", values.room);
    if(values.room !== undefined && values.room !== ""){
      props.getMeta(values.room)
    }
  }, [props.location.search]);

  useEffect(() => {
    form.setFieldsValue({
      roomname: props.room.name,
      roomid: props.room.id
    });
  }, [props.room]);
 

  return (
    <div className="wrap-login">
        <Content className="login-form">
            <Typography>
            <p style={{"width": "500px", marginBottom:"70px"}}>
              <img src={logo} style={{"float": "left", "width": "80px", marginRight:"30px"}} />
              <Title style={{width: "400px",  display:"block", paddingTop:"16px"}} level={1}>Watch2Gether</Title>
            </p> 
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
            
            <Form
              {...layout}
              name="basic"
              form={form}
              initialValues={{ anonymous: true }}
              onFinish={onFinish}
              onFinishFailed={onFinishFailed}
            >
               {/* <Form.Item
                label="Room id"
                name="roomid"
               >
                <Input size="large" />
              </Form.Item> */}

              <Form.Item
                label="Room Name"
                name="roomname"
                rules={[{ required: true, message: 'Please input your room Name!' }, { min: 4, message: 'Room name must be minimum 4 characters.' }]}
              >
                <Input size="large" />
              </Form.Item>

              <Form.Item
                label="Username"
                name="username"
                rules={[{ required: true, message: 'Please input your username!' }, { min: 4, message: 'Username must be minimum 4 characters.' }]}
                style={{ "marginTop": "20px" }}
              >
                <Input size="large" />
              </Form.Item>

              <Form.Item {...tailLayout} name="anonymous" valuePropName="checked">
                <Checkbox>be anonymous </Checkbox>
              </Form.Item>

              <Form.Item {...tailLayout}>
                <Button size="large" type="primary" htmlType="submit" style={{ "width": "100%", marginTop: "20px" }}>
                  Login
                </Button>

                <Divider/>

                {/* <Button size="large" type="primary" style={{ "width": "100%", marginTop: "0px", backgroundColor: "#7289da", border: "none" }}>
                  Login with Discord 
                </Button> */}
              </Form.Item>
            </Form>
            <PageFooter style={{ textAlign: 'center', position: "absolute", bottom:"0px", left:"0px", width:"560px", height:"50px", padding: "15px 50px 28px 50px" }}/>
        </Content>
      </div>
    
  );
}

const mapStateToProps = (state) => {

  return ({ error: state.room.error, room: state.room })
}
export default withRouter(connect(mapStateToProps, { join, leave, clearError, getMeta })(Home))


