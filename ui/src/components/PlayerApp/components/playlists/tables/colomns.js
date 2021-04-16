import { sortableHandle } from 'react-sortable-hoc';
import { MenuOutlined } from '@ant-design/icons';

import { Popconfirm } from 'antd';
import {DeleteOutlined} from '@ant-design/icons';

const DragHandle = sortableHandle(() => <MenuOutlined style={{ cursor: 'grab', color: '#999' }} />);

const cols = [
        {
        title: 'Video Address',
        dataIndex: 'url',
        editable: true,
        }
]

export const EditableCols = (handleDelete) => {
    return [
        {
        title: '',
        dataIndex: 'Delete',
        width: 80,
        render: (_, record) => (
            <Popconfirm title="Sure to delete?" onConfirm={() => handleDelete(record.key)}>
                <DeleteOutlined />
            </Popconfirm>
            ),
        },
        ...cols
    ];
}

export const SortableCols = () => {
    return [
    {
      title: 'Sort',
      dataIndex: 'sort',
      width: 80,
      className: 'drag-visible',
      render: () => <DragHandle />,
    },
    ...cols
  ];
}