import React from "react"
import { Avatar } from 'antd';

const UserIcon = ({ icon, username }) => {
  if (icon !== undefined && icon !== '') {
    return (<Avatar shape="circle" size={38} src={icon} />);
  }
  return (
    <Avatar style={{ verticalAlign: 'middle', marginRight: '14px' }} shape="circle" size={38}>
      {username.substring(0, 2)}
    </Avatar>
  );
};
export default UserIcon;
