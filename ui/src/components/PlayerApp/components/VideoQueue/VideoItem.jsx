import React from 'react';
import { List, Typography, Skeleton } from 'antd';

import VideoThumbnail  from './VideoThumbnail';

const { Title } = Typography;

export const VideoItem = ({
  video, children, playing, loading,
}) => (
  <List.Item>
    <table className="">
      <tbody>
        <tr>
          {video.thumbnail&&
            <td style={{ width: "130px" }}>
              {<VideoThumbnail url={video.thumbnail} user={video.user} />}
            </td>
          }
          <td style={{
            height: '50px', 
            overflowY: 'hidden', 
            display: 'inline-block',
            maxWidth: video.thumbnail?'250px':"300px",
            margin: '0px',
            padding: '0px',
            overflow: 'hidden',
          }}
          >
            {playing ? 'Currently Playing' : null}
            {!loading
              ? (
                <div>
                  <Title level={5} style={{ fontSize: '14px', margin:"5px" }} className="eclipseText">
                    {video.title}
                  </Title>
                  Added by:
                  {' '}
                  {video.user}
                </div>
              )
              : <Skeleton size="small" active /> }
          </td>
        </tr>
      </tbody>
    </table>
    {children !== undefined && !loading
      ? (
        <div
          className="videoQueueItem"
          style={{
            position: 'absolute',
            width: '620px',
          }}
        >
          <div style={{
            float: 'right', background: '#141414', padding: '15px 20px 15px 5px', width: '280px',
          }}
          >
            {children}
          </div>
        </div>
      )
      : null}
  </List.Item>
);
