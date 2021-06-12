import React from "react"
import { sortableHandle } from 'react-sortable-hoc';
import { MenuOutlined, DeleteOutlined } from '@ant-design/icons';

import { Popconfirm, Skeleton } from 'antd';

import  VideoThumbnail  from '../../VideoQueue/VideoThumbnail';

const DragHandle = sortableHandle(() => <MenuOutlined style={{ cursor: 'grab', color: '#999' }} />);
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

export const EditableCols = (handleDelete) => [
  {
    title: 'Delete',
    dataIndex: 'Delete',
    width: 80,
    render: (_, record) => (
      <Popconfirm title="Sure to delete?" onConfirm={() => handleDelete(record.key)}>
        <DeleteOutlined />
      </Popconfirm>
    ),
  },

  ...cols,
];

export const SortableCols = () => [
  {
    title: 'Sort',
    dataIndex: 'sort',
    width: 80,
    className: 'drag-visible',
    render: () => <DragHandle />,
  },
  ...cols,
];
