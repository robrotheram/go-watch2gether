
import { Card, List, Avatar, Button } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import {connect} from 'react-redux'
import {updateHost} from '../store/room/room.actions'
function UserList(props){
    const {users, isHost, host } = props;

    const listActions = (item) => {
      let actions = [] 
      if(isHost  && host!==item.name){
        actions.push(<Button type="link"  key="list-loadmore-edit" onClick={() => updateHost(item.name)}>Promote To Host</Button>)
      }
      // if (isHost) {
      //   actions.push( <a key="list-loadmore-more">Kick</a>);
      // }
      return actions;
    }
    return (
      <Card type="inner" title="Users List:" className="list">
        <div className="container .sc2">
          <List
            size="small"
            itemLayout="horizontal"
            dataSource={users}
            renderItem={item => (
              <List.Item actions={listActions(item)}>
                <div><Avatar icon={<UserOutlined />}  style={{"marginRight": "10px"}}/>{item.name}</div>
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

