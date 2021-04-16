import { Avatar } from 'antd';
const UserIcon = ({user}) => {
    if (user.icon !== undefined && user.icon !== ""){
        return (<Avatar shape="circle" size={38} src={user.icon} />)
    }
    return (
        <Avatar style={{verticalAlign: 'middle', marginRight:"14px"}} shape="circle" size={38} >
            {user.username.substring(0, 2)}
        </Avatar>
    )
}
export default UserIcon