import './App.less';
import { useEffect } from "react"
import { Layout, Button, Divider, Card, Typography } from 'antd';
import { Form, Input } from 'antd';
import { Alert } from 'antd';
import {connect} from 'react-redux'
import {join, leave, clearError} from './store/room/room.actions'
import { PageFooter } from './components/PageFooter';
import { withRouter } from "react-router";
import queryString from 'query-string'

const { Content } = Layout;
const { Title, Paragraph, Text } = Typography;


function Home(props) {
  const roomName = queryString.parse(props.location.search).room;
  const layout = {
    labelCol: { span: 4 },
    wrapperCol: { span: 20 },
  };
  const tailLayout = {
    wrapperCol: { offset: 0, span: 4 },
  };
  

  const onFinish = values => {
    console.log("Valuse", values)
    props.join(values.roomname, values.username)     
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
  }, [props.location.search]);
  

  return (
    <Layout className="dark-theme" style={{"backgroundColor" : "white"}}>
    <Content style={{ padding: '88px 100px', position: "fixed", top:"65px", left:"0px", right:"0px", bottom:"50px", overflow:"auto" }}>
    <Card style={{"padding":"50px"}}>
      <Typography>
        <Title>Watch2Gether</Title>
        <Paragraph>
          The hot new video synchronization platform used to watch videos in realtime with friends!
        </Paragraph>
        <Paragraph>
          Enjoy content from YouTube, Vimeo, Dailymotion and SoundCloud as well as your own media files
          <br/>
          <Text >
            It the easy solution, which lets you create your own room to add multiple users from different parts of the world at the same time.
          </Text>
        </Paragraph>
      </Typography>

    <Divider/>

    {props.error !== "" ? (
      <Alert
        message="Error"
        description={props.error}
        type="error"
        showIcon
        closable
        afterClose={handleClose}
        style={{"marginBottom":"20px"}}
      />
    ) : null}
    
    <Form
      {...layout}
      name="basic"
      initialValues={{ roomname: roomName}}
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
    >
      <Form.Item
        label="Room Name"
        name="roomname"
        rules={[{ required: true, message: 'Please input your room Name!' }]}
      >
        <Input size="large"/>
      </Form.Item>

      <Form.Item
        label="Username"
        name="username"
        rules={[{ required: true, message: 'Please input your username!' }]}
        style={{"marginTop": "20px"}}
      >
        <Input size="large"/>
      </Form.Item>

      <Form.Item {...tailLayout}>
        <Button size="large" type="primary" htmlType="submit" style={{"width":"100%", marginTop: "20px"}}>
          Submit
        </Button>
      </Form.Item>
    </Form>


    </Card>

    </Content>
    <PageFooter/>
    </Layout>
  );
}

const mapStateToProps  = (state) =>{
  
  return ({error:state.room.error})
} 
export default withRouter(connect(mapStateToProps, {join, leave, clearError})(Home))


