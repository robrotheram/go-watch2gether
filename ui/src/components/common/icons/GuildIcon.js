import { Avatar } from 'antd';
const style = {height:"100%",width:"100%", paddingRight:"14px", marginRight:"12px", verticalAlign: 'middle'}
const GuildIcon = ({guild}) => {
    if (guild.icon !== undefined && guild.icon !== ""){
        return (<img style={style}  src= {"https://cdn.discordapp.com/icons/"+guild.id+"/"+guild.icon+".png"} />)
    }
    return (
        <Avatar style={{verticalAlign: 'middle', marginRight:"14px"}} shape="square" size={55} >
            {guild.name.substring(0, 2)}
        </Avatar>
    )
}
export default GuildIcon