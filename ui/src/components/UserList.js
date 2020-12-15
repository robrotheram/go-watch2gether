
import { Card, List, Button, Progress } from 'antd';
import { Row, Col } from 'antd';
import {connect} from 'react-redux'
import {updateHost} from '../store/room/room.actions'
function UserList(props){
    const {users, isHost, host, user } = props;

    const listActions = (item) => {
      let actions = [] 
      if(isHost){
        actions.push(<Button type="link" disabled={isHost  && host===item.name}  key="list-loadmore-edit" onClick={() => updateHost(item.name)}>Promote To Host</Button>)
      }
      // if (isHost) {
      //   actions.push( <a key="list-loadmore-more">Kick</a>);
      // }
      return actions;
    }
    return (
      <Card type="inner" title="Users Progress" className="list">
        <div className="container .sc2 userlist">
          {/* {JSON.stringify(user)} */}
          <List
            size="small"
            itemLayout="horizontal"
            dataSource={users}
            renderItem={item => (
              <List.Item actions={listActions(item)} className={item.name === host ? "userListActive" : null}>
                  <Row style={{"width":"100%", "padding":"5px"}}>
                    <Col flex="150px" style={{"textAlign":"right", "paddingRight":"10px"}}>
                      {item.name === host ? "Host: "+item.name : item.name}  
                    </Col>
                    <Col flex="auto" >
                      <div style={{"display":"inline-block", "width":"100%"}}><Progress percent={(item.seek)*100} showInfo={false}size="small"/></div>
                    </Col>
                </Row>
              </List.Item>
            )}
          />
        </div>
      </Card>
    )
}

const mapStateToProps  = (state) =>{
  
  return state.room
} 
export default connect(mapStateToProps, {updateHost})(UserList)

