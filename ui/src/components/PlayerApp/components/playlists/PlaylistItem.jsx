import React from "react"
import { List, Typography, Skeleton } from 'antd';

import  VideoThumbnail  from '../VideoQueue/VideoThumbnail';

const { Title } = Typography;

export const PlaylistItem = ({
  video, children, playlist, playing, loading, click,
}) => {
  let url = '';
  let user = '';
  if (video !== undefined) {
    if (video.url !== undefined) {
      url = video.url;
      user = video.user;
    }
  }

  return (
    <List.Item style={{ cursor: 'pointer' }} onClick={() => click()} className="playlistItemContainer">
      <table>
        <tbody>
          <tr>
            <td style={{ width: '130px' }}>
              {loading || url === '' || video === undefined
                ? <Skeleton.Image style={{ height: '70px', padding: '10px' }} />
                : <VideoThumbnail url={url} user={user} />}
            </td>
            <td style={{
              height: '70px',
              overflowY: 'hidden',
              padding: '0px 10px',
              width: '240px',
              display: 'inline-block',
            }}
            >
              {playing ? 'Currently Playing' : null}
              {!loading
                ? (
                  <div>
                    <Title level={5} style={{ fontSize: '14px' }} className="eclipseText">
                      {playlist.name}
                    </Title>
                    <p style={{ marginBottom: '0px' }}>
                      Last modified by:
                      {playlist.username}
                    </p>
                    <p>
                      Length:
                      {playlist.videos.length}
                    </p>
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
            className="playlistItem"
            style={{
              position: 'absolute',
              width: '910px',
              marginTop: '-5px',
            }}
          >
            <div style={{
              float: 'right',

              padding: '20px 140px 20px 12px',
            }}
            >
              {children}
            </div>
          </div>
        )
        : null}
    </List.Item>
  );
};
