
import { List, Button, Progress } from 'antd';
import { Row, Col } from 'antd';
import {useSelector, useDispatch} from 'react-redux'
import { SyncOutlined } from '@ant-design/icons';
import { seekToUser } from '../../../store/video/video.actions';

function UserList(){
  const dispatch = useDispatch()

  const watchers = useSelector(state => state.room.watchers);
  const host = useSelector(state => state.room.host);

    const secondsToDate = (seconds) => {
      let time =  new Date(seconds * 1000).toISOString().substr(11, 8)
      var res = time.substring(0,2)
      if (res === "00") {
        return time.substring(3,time.length)
      }
      return time
    }
    //console.log("room",props.room.watchers)
    return (
      // <Card type="inner" title="Users Progress" className="list">
      //   <div className="container .sc2 userlist">
          /* {JSON.stringify(watchers)} */
          <List
            size="small"
            itemLayout="horizontal"
            dataSource={watchers}
            renderItem={item => (
              <List.Item className={item.id === host ? "userListActive" : null}>
                  <Row style={{"width":"100%", "padding":"5px"}}>
                    <Col flex="100px" style={{"textAlign":"left", "paddingRight":"10px"}}>
                      {item.username}  
                    </Col>
                    <Col flex="auto" >
                      <div style={{"display":"inline-block", "width":"100%", paddingRight:"10px"}}>
                        <Progress percent={(item.seek.progress_percent)*100} showInfo={false}size="small"/>
                      </div>
                    </Col>
                    <Col>
                       {secondsToDate(item.seek.progress_seconds)}
                       <Button icon={<SyncOutlined />} onClick={() => {dispatch(seekToUser(item.seek))}} style={{marginLeft: "10px"}}/>
                    </Col>
                </Row>
              </List.Item>
            )}
          />
      //   </div>
      // </Card>
    )
}

export default UserList

