

import { Card, List, Space, Button,Typography  } from 'antd';
import { ArrowUpOutlined, ArrowDownOutlined, DeleteOutlined } from '@ant-design/icons';

import {VideoThumbnail} from "./VideoThumbnail"
import {connect} from 'react-redux'
import {updateQueue} from '../../store/room/room.actions'
const { Title, Paragraph, Text, Link } = Typography;

export function VideoListComponent (props) {

    const { queue } = props

  const deleteVideo = (item) => {
    let videoList = [...queue];
    let i = videoList.indexOf(item);
    videoList.splice(i, 1);
    updateQueue(videoList)
  }

  const voteUp = (item) => {
    let videoList = [...queue];
    let i = videoList.indexOf(item);
    var z = videoList[i-1];
    videoList[i-1] = videoList[i];
    videoList[i] = z;
    updateQueue(videoList)
  }

  const voteDown = (item) => {
    let videoList = [...queue]; 
    let i = videoList.indexOf(item);
    var z = videoList[i+1];
    videoList[i+1] = videoList[i];
    videoList[i] = z;
    updateQueue(videoList)
  }

  
    return(
        <Card type="inner" className="list" title={<Paragraph>There are <Text strong>{queue.length}</Text> videos in the queue</Paragraph>}>
            <div className="videoQueue">
            <List
                size="small"
                itemLayout="horizontal"
                dataSource={queue}
                renderItem={item => (
                    <List.Item className={queue.indexOf(item) === 0 ? "itemPlaying":""}>
                    <table className="videoQueueItem">
                        <tbody>
                        <tr>
                            <td style={{"width":"130px"}}> 
                            <VideoThumbnail url={item.url} user={item.user}/>
                            </td>
                            <td style={{"padding":"0px 10px", "maxWidth":"250px"}}> 
                            <Title level={5}  style={{fontSize:"14px"}} className="eclipseText">
                                {item.url}
                            </Title>
                            Added by: {item.user}
                            </td>
                            <td style={{"width":"112px"}}>
                            <Space>
                                { queue.indexOf(item) !== 0 ? <Button onClick={() => voteUp(item)} icon={<ArrowUpOutlined />} /> : null}
                                { queue.indexOf(item) !== queue.length - 1? <Button onClick={() => voteDown(item)} icon={<ArrowDownOutlined />} /> : null}
                                <Button icon={<DeleteOutlined onClick={() => deleteVideo(item)}/>} />
                            </Space>
                            </td>
                        </tr>
                        </tbody>
                    </table>
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
export const VideoList = connect(mapStateToProps, {updateQueue})(VideoListComponent)
  
  