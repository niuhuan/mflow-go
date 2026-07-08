// 模板工程直接引用 assets 中的 XML 文件（方便将来替换），而非内联字符串。
import emptyXml from '../assets/templates/empty.xml?raw';
import singleUserXml from '../assets/templates/single-user.xml?raw';
import multiAccountXml from '../assets/templates/mutil-account.xml?raw';

export { emptyXml, singleUserXml, multiAccountXml };

export const templates = [
  { key: 'empty', name: '空白工程', xml: emptyXml },
  { key: 'single', name: '单账号模板', xml: singleUserXml },
  { key: 'multi', name: '多账号模板', xml: multiAccountXml },
];
