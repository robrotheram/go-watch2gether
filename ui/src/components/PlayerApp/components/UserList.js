
import { Card, List, Button, Progress } from 'antd';
import { Row, Col } from 'antd';
import {connect} from 'react-redux'
import {updateHost} from '../../../store/room/room.actions'
// const watchers = [
//   {name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"},{name:"test"}
// ]

function UserList(props){
    const {watchers, host } = props.room;
    const {isHost} = props.user;

    const listActions = (item) => {
      let actions = [] 
      if(isHost){
        actions.push(
          <Button 
            type="link" 
            disabled={isHost  && host===item.id}  
            key="list-loadmore-edit" 
            onClick={() => updateHost(item.name)}>
              {isHost  && host===item.id ?"You are the host": "Promote To Host"}
          </Button>
        )
      }
      return actions;
    }
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
                      {item.name}  
                    </Col>
                    <Col flex="auto" >
                      <div style={{"display":"inline-block", "width":"100%"}}><Progress percent={(item.seek)*100} showInfo={false}size="small"/></div>
                    </Col>
                </Row>
              </List.Item>
            )}
          />
      //   </div>
      // </Card>
    )
}

const mapStateToProps  = (state) =>{
  return state
} 
export default connect(mapStateToProps, {updateHost})(UserList)

