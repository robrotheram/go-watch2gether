import React from "react"
import { Avatar } from 'antd';

const style = {
  height: '100%', width: '100%', paddingRight: '14px', marginRight: '12px', verticalAlign: 'middle', borderRadius: '50%',
};
const GuildIcon = ({ guild }) => {
  if (guild.icon !== undefined && guild.icon !== '') {
    return (<img style={style} alt={guild.name} src={`https://cdn.discordapp.com/icons/${guild.id}/${guild.icon}.png`} />);
  }
  return (
    <Avatar style={{ verticalAlign: 'middle', marginRight: '14px', borderRadius: '50%' }} size={55}>
      {guild.name.substring(0, 2)}
    </Avatar>
  );
};
export default GuildIcon;
