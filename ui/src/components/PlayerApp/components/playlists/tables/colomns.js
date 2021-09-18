import React from "react"
import { DeleteOutlined, ArrowUpOutlined, ArrowDownOutlined } from '@ant-design/icons';

import { Popconfirm, Skeleton, Button } from 'antd';

import  VideoThumbnail  from '../../VideoQueue/VideoThumbnail';

/*eslint react/display-name: "off"*/
export const cols = [

  {
    title: 'Icon',
    dataIndex: 'url',
    width: 100,
    render: (_, record) => {
      if (record.url === undefined) {
        return <Skeleton.Image style={{ height: '70px', padding: '10px' }} />;
      }
      return <div style={{ width: 100 }}><VideoThumbnail url={record.url} /></div>;
    },
  },
  {
    title: 'Video Address',
    dataIndex: 'url',
    editable: true,
  },
];

export const EditableCols = (handleDelete, handleMove) => [
  {
    title: 'Actions',
    dataIndex: 'Actions',
    width: 80,
    render:( text, record, pos) => (
      <table style={{width:"100%", border:"0"}}>
        <tr>
          <td>
            <span>
              <Button type="default" icon={<ArrowUpOutlined />}  onClick={() => handleMove(pos, (pos-1))}/>
              <Button type="default" icon={<ArrowDownOutlined />} onClick={() => handleMove(pos, (pos+1))}/>
            </span>
          </td>
          <td style={{padding:"5px"}}>
            <Popconfirm title="Sure to delete?" onConfirm={() => handleDelete(record.key)}>
            <Button type="default" icon={<DeleteOutlined />}/>
            </Popconfirm>
          </td>
        </tr>
      </table>
    ),
  },
  ...cols
];