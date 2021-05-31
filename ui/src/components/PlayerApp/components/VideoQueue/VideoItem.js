
import {List,Typography  } from 'antd';
import {VideoThumbnail} from "./VideoThumbnail"
import { Skeleton } from 'antd';
const { Title } = Typography;


export const VideoItem = ({video, children, playing, loading}) => {
    return (
        <List.Item>
                <table className="">
                    <tbody>
                    <tr>
                        <td style={{"width":"130px"}}> 
                        {loading ? 
                            <Skeleton.Image style={{"height":"60px", "padding": "10px"}} /> 
                            :
                            <VideoThumbnail url={video.url} user={video.user}/>
                        }
                        </td>
                        <td style={{"height":"50px", "overflowY":"hidden", "padding":"0px 10px", "width":"240px", "display":"inline-block"}}> 
                        {playing ?"Currently Playing":null}
                        {!loading ? 
                        <div>
                            <Title level={5}  style={{fontSize:"14px"}} className="eclipseText">
                                {video.title}
                            </Title>
                            Added by: {video.user}
                        </div>
                        : <Skeleton size={"small"} active/> }
                        </td>
                    </tr>
                    </tbody>
                </table>
                {children !== undefined && !loading ? 
                <div className="videoQueueItem" style={{
                    position:"absolute",
                    width: "620px"
                }}>
                    <div style={{float: "right", background: "#141414", padding: "15px 20px 15px 5px", width:"260px"}}>
                        {children}
                    </div>
                </div>
                :null}
        </List.Item>
    )
}