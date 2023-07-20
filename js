遍历目录，找到最后一个含有文件的目录，然后做操作，js代码：
const fs = require('fs');
const util = require('util');
const path = require('path');
const { DirectoryLoader, TextLoader } = require('some-library'); // Replace 'some-library' with the actual library you are using

const readdir = util.promisify(fs.readdir);
const stat = util.promisify(fs.stat);


async function findTxtFilesDirectory(chatid) {
  async function traverseDirectory(dir) {
    const files = await readdir(dir);

    for (const file of files) {
      const filePath = path.join(dir, file);
      const fileStat = await stat(filePath);

      if (fileStat.isDirectory()) {
        // Recursively traverse the subdirectories
        const subdirectory = await traverseDirectory(filePath);
        if (subdirectory) {
          return subdirectory; // Found a subdirectory with .txt files
        }
      } else if (file.endsWith('.txt')) {
        // If a .txt file is found, return the current directory
        return dir;
      }
    }

    return null; // No .txt files found in this directory or its subdirectories
  }

  // Start traversing from the specified chatid directory
  const rootDir = path.join('docs', chatid);
  const txtFilesDir = await traverseDirectory(rootDir);
  return txtFilesDir;
}

// Usage
const chatid = 'parent-directory';
ingestData(chatid)
  .then((result) => {
    console.log('Directory containing .txt files:', result);
  })
  .catch((error) => {
    console.error('Error:', error);
  });




