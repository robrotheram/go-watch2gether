

import { Card, List, Space, Button,Typography  } from 'antd';
import { ArrowUpOutlined, ArrowDownOutlined, DeleteOutlined } from '@ant-design/icons';

import {connect} from 'react-redux'
import {updateQueue} from '../../../../store/room/room.actions'
import {VideoItem} from './VideoItem'
const { Paragraph, Text } = Typography;

export function VideoListComponent (props) {

  const { queue } = props.room
  const current_video = props.video

  const skipTo = (item) => {
    let videoList = [...queue];
    let i = videoList.indexOf(item);
    videoList.splice(i, 1);
    videoList.unshift(item)
    updateQueue(videoList)

  }

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
      <div style={{"height":"100%"}}>
        {current_video !== undefined && current_video.url !== "" ? 
         <VideoItem video={current_video} playing/>
        : null}
        <Card type="inner" 
          className="list video"
          title={<Paragraph>There are <Text strong>{queue.length}</Text> videos in the queue</Paragraph>}
          style={ current_video !== undefined && current_video.url !== "" ? {"height": "calc( 100% - 156px )"} : {"height": "calc( 100% - 52px )"}}
         >
            <div className="videoQueue">
            <List
                size="small"
                itemLayout="horizontal"
                dataSource={queue}
                renderItem={item => (
                      <VideoItem video={item} loading={item.loading}>
                        <Space>
                                 { queue.indexOf(item) !== 0 ? <Button onClick={() => voteUp(item)} icon={<ArrowUpOutlined />} /> : null}
                                 { queue.indexOf(item) !== queue.length - 1? <Button onClick={() => voteDown(item)} icon={<ArrowDownOutlined />} /> : null}
                                 <Button icon={<DeleteOutlined onClick={() => deleteVideo(item)}/>} />
                                 { props.isHost  && queue.indexOf(item) > 0  ? <Button onClick={() => skipTo(item) }>Move To Top</Button> : null }
                        </Space>
                    </VideoItem>
                    // <List.Item>
                    // <table className="videoQueueItem">
                    //     <tbody>
                    //     <tr>
                    //         <td style={{"width":"130px"}}> 
                    //         <VideoThumbnail url={item.url} user={item.user}/>
                    //         </td>
                    //         <td style={{"padding":"0px 10px", "maxWidth":"250px"}}> 
                    //         <Title level={5}  style={{fontSize:"14px"}} className="eclipseText">
                    //             {item.url}
                    //         </Title>
                    //         Added by: {item.user}
                    //         </td>
                    //         <td style={{"width":"250px"}}>
                    //         
                    //         </td>
                    //     </tr>
                    //     </tbody>
                    // </table>
                    // </List.Item>
                )}
                />
            </div>
        </Card>
        </div>
    )
}
const mapStateToProps  = (state) =>{
    
    return state
  } 
export const VideoList = connect(mapStateToProps, {updateQueue})(VideoListComponent)
  
  