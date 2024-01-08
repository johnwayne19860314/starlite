// import React, { useState } from 'react';
// import { Table, Collapse } from 'antd';
// import 'antd/dist/antd.css';

// const { Panel } = Collapse;

// interface DataSourceItem {
//   key: string;
//   name: string;
//   age: number;
//   address: string;
// }

// const dataSource: DataSourceItem[] = [
//   {
//     key: '1',
//     name: 'John Doe',
//     age: 25,
//     address: '123 Main St',
//   },
//   {
//     key: '2',
//     name: 'Jane Doe',
//     age: 30,
//     address: '456 Park Ave',
//   },
// ];

// const columns = [
//   {
//     title: 'Name',
//     dataIndex: 'name',
//     key: 'name',
//   },
//   {
//     title: 'Age',
//     dataIndex: 'age',
//     key: 'age',
//   },
//   {
//     title: 'Address',
//     dataIndex: 'address',
//     key: 'address',
//   },
// ];

// const CollapseTable: React.FC = () => {
//   const [expandedRowKeys, setExpandedRowKeys] = useState<string | string[]>([]);

//   const handleExpand = (expandedRows: string | string[]) => {
//     setExpandedRowKeys(expandedRows);
//   };

//   const expandedRowRender = (record: DataSourceItem) => (
//     <p style={{ margin: 0 }}>{record.name} details go here.</p>
//   );

//   return (
//     <Collapse
//       accordion
//       onChange={handleExpand}
//       activeKey={expandedRowKeys}
//     >
//       {dataSource.map((record) => (
//         <Panel
//           header={record.name}
//           key={record.key}
//         >
//           <Table
//             dataSource={[record]}
//             columns={columns}
//             pagination={false}
//             expandedRowRender={expandedRowRender}
//             rowKey="key"
//           />
//         </Panel>
//       ))}
//     </Collapse>
//   );
// };

// export default CollapseTable;
