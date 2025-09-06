<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import { FileJson, FileSpreadsheet, FileCode, Download, Copy, Check } from 'lucide-svelte';
  import type { Storyboard } from '../../types/storyboard';

  interface Props {
    storyboard?: Storyboard;
    closeModal?: () => void;
  }

  let { storyboard, closeModal = () => {} }: Props = $props();

  // State for managing generated downloads
  let downloads: {
    json?: { url: string; filename: string; data: string };
    csv?: { url: string; filename: string; data: string };
    xml?: { url: string; filename: string; data: string };
  } = $state({});

  let copiedStates = $state({
    json: false,
    csv: false,
    xml: false
  });

  // For screen reader announcements
  let announceText = $state('');

  function safeFilename(name: string, ext: string) {
    const base = name ? name.replace(/[^a-z0-9-_]+/gi, '_').replace(/_{2,}/g,'_').replace(/^_|_$/g,'') : 'storyboard';
    return `${base}-storyboard.${ext}`;
  }

  function createDownloadLink(data: string, filename: string, type: string) {
    try {
      const blob = new Blob([data], { type });
      const url = URL.createObjectURL(blob);
      return { url, filename, data };
    } catch (error) {
      console.error('Failed to create download link:', error);
      return null;
    }
  }

  function generateJSON() {
    if (!storyboard) return;
    const {
         facilitatorCode: _facilitatorCode,
         joinCode: _joinCode,
         facilitators: _facilitators,
         users: _users,
         owner_id: _owner_id,
         teamId: _team_id,
         teamName: _team_name,
         ...safe
       } = storyboard;

    const jsonData = JSON.stringify(safe, null, 2);
    const filename = safeFilename(storyboard.name, 'json');
    const result = createDownloadLink(jsonData, filename, 'application/json');
    if (result) {
      downloads.json = result;
    }
  }

  function generateCSV() {
    if (!storyboard) return;
    const rows: string[] = [];
    rows.push('Goal,Column,Story,Points,Closed,Color,Link');
    storyboard.goals.forEach(g => {
      g.columns.forEach(c => {
        c.stories.forEach(s => {
          const safe = (val: string | number | boolean) => `"${String(val).replace(/"/g, '""')}"`;
          rows.push([
            safe(g.name),
            safe(c.name),
            safe(s.name),
            s.points ?? 0,
            s.closed,
            safe(s.color),
            safe(s.link || '')
          ].join(','));
        });
      });
    });
    const csvContent = '\uFEFF' + rows.join('\n');
    const filename = safeFilename(storyboard.name, 'csv');
    const result = createDownloadLink(csvContent, filename, 'text/csv;charset=utf-8;');
    if (result) {
      downloads.csv = result;
    }
  }

  function generateXML() {
    if (!storyboard) return;
    const esc = (v: string) => v.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
    let xml = '<?xml version="1.0" encoding="UTF-8"?>\n';
    xml += `<storyboard id="${esc(storyboard.id)}" name="${esc(storyboard.name)}">`;
    storyboard.goals.forEach(g => {
      xml += `\n  <goal id="${esc(g.id)}" name="${esc(g.name)}">`;
      g.columns.forEach(c => {
        xml += `\n    <column id="${esc(c.id)}" name="${esc(c.name)}">`;
        c.stories.forEach(s => {
          xml += `\n      <story id="${esc(s.id)}" name="${esc(s.name)}" points="${s.points}" closed="${s.closed}" color="${esc(s.color)}">`;
          if (s.link) xml += `\n        <link>${esc(s.link)}</link>`;
          if (s.comments && s.comments.length) {
            xml += '\n        <comments>';
            s.comments.forEach(cm => {
              xml += `\n          <comment id="${esc(cm.id)}">${esc(cm.comment)}</comment>`;
            });
            xml += '\n        </comments>';
          }
          xml += '\n      </story>';
        });
        xml += '\n    </column>';
      });
      xml += '\n  </goal>';
    });
    xml += '\n</storyboard>';
    const filename = safeFilename(storyboard.name, 'xml');
    const result = createDownloadLink(xml, filename, 'application/xml');
    if (result) {
      downloads.xml = result;
    }
  }

  async function copyToClipboard(data: string, format: 'json' | 'csv' | 'xml') {
    try {
      await navigator.clipboard.writeText(data);
      copiedStates[format] = true;
      announceText = `${format.toUpperCase()} data copied to clipboard`;
      
      setTimeout(() => {
        copiedStates[format] = false;
        announceText = '';
      }, 2000);
    } catch (error) {
      console.error('Failed to copy to clipboard:', error);
      // Fallback: select text in a temporary textarea
      const textarea = document.createElement('textarea');
      textarea.value = data;
      textarea.style.position = 'fixed';
      textarea.style.opacity = '0';
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand('copy');
      document.body.removeChild(textarea);
      
      // Still show success feedback for fallback
      copiedStates[format] = true;
      announceText = `${format.toUpperCase()} data copied to clipboard`;
      
      setTimeout(() => {
        copiedStates[format] = false;
        announceText = '';
      }, 2000);
    }
  }

  // Cleanup URLs when component is destroyed
  function cleanup() {
    Object.values(downloads).forEach(download => {
      if (download?.url) {
        URL.revokeObjectURL(download.url);
      }
    });
  }

  // Generate all downloads on component mount
  $effect(() => {
    if (storyboard) {
      generateJSON();
      generateCSV();
      generateXML();
    }
    
    return cleanup;
  });
</script>

<Modal closeModal={closeModal} ariaLabel="Storyboard Export Options" widthClasses="md:w-4/5 lg:w-3/4 xl:w-2/3">
  <div>
    <!-- Screen reader announcements -->
    <div aria-live="polite" aria-atomic="true" class="sr-only">
      {announceText}
    </div>
    
    <!-- Skip link for keyboard users -->
    <a href="#export-actions" class="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 bg-blue-600 text-white px-4 py-2 rounded-lg z-50">
      Skip to export actions
    </a>
    
    <h2 class="text-2xl font-rajdhani mb-6 text-gray-900 dark:text-white">Export Storyboard</h2>
    <div class="space-y-6" data-testid="storyboard-export-options" role="main" id="export-actions">
      
      <!-- JSON Export -->
      <div class="border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 rounded-lg p-6 shadow-sm dark:shadow-gray-900/20">
        <div class="flex items-center justify-between mb-4">
          <div>
            <div class="font-bold flex items-center space-x-2 text-gray-900 dark:text-white">
              <FileJson class="w-5 h-5 text-green-600 dark:text-green-400" />
              <span>JSON Export</span>
            </div>
            <div class="text-gray-600 dark:text-gray-300 mt-1">Full raw data dump with complete structure</div>
          </div>
        </div>
        
        {#if downloads.json}
          <div class="flex flex-wrap gap-3" role="group" aria-label="JSON export actions">
            <SolidButton 
              color="green"
              href={downloads.json.url}
              options={{ 
                download: downloads.json.filename,
                'aria-label': `Download JSON file: ${downloads.json.filename}`
              }}
            >
              <Download class="w-4 h-4 mr-2" aria-hidden="true" />
              Download {downloads.json.filename}
            </SolidButton>
            
            <SolidButton
              color="blue"
              onClick={() => copyToClipboard(downloads.json!.data, 'json')}
              options={{ 'aria-label': 'Copy JSON data to clipboard' }}
            >
              {#if copiedStates.json}
                <Check class="w-4 h-4 mr-2" aria-hidden="true" />
                Copied!
              {:else}
                <Copy class="w-4 h-4 mr-2" aria-hidden="true" />
                Copy Data
              {/if}
            </SolidButton>
          </div>
        {/if}
      </div>

      <!-- CSV Export -->
      <div class="border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 rounded-lg p-6 shadow-sm dark:shadow-gray-900/20">
        <div class="flex items-center justify-between mb-4">
          <div>
            <div class="font-bold flex items-center space-x-2 text-gray-900 dark:text-white">
              <FileSpreadsheet class="w-5 h-5 text-orange-600 dark:text-orange-400" />
              <span>CSV Export</span>
            </div>
            <div class="text-gray-600 dark:text-gray-300 mt-1">Tabular format for spreadsheets (Excel, Google Sheets)</div>
          </div>
        </div>
        
        {#if downloads.csv}
          <div class="flex flex-wrap gap-3" role="group" aria-label="CSV export actions">
            <SolidButton 
              color="green"
              href={downloads.csv.url}
              options={{ 
                download: downloads.csv.filename,
                'aria-label': `Download CSV file: ${downloads.csv.filename}`
              }}
            >
              <Download class="w-4 h-4 mr-2" aria-hidden="true" />
              Download {downloads.csv.filename}
            </SolidButton>
            
            <SolidButton
              color="blue"
              onClick={() => copyToClipboard(downloads.csv!.data, 'csv')}
              options={{ 'aria-label': 'Copy CSV data to clipboard' }}
            >
              {#if copiedStates.csv}
                <Check class="w-4 h-4 mr-2" aria-hidden="true" />
                Copied!
              {:else}
                <Copy class="w-4 h-4 mr-2" aria-hidden="true" />
                Copy Data
              {/if}
            </SolidButton>
          </div>
        {/if}
      </div>

      <!-- XML Export -->
      <div class="border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 rounded-lg p-6 shadow-sm dark:shadow-gray-900/20">
        <div class="flex items-center justify-between mb-4">
          <div>
            <div class="font-bold flex items-center space-x-2 text-gray-900 dark:text-white">
              <FileCode class="w-5 h-5 text-purple-600 dark:text-purple-400" />
              <span>XML Export</span>
            </div>
            <div class="text-gray-600 dark:text-gray-300 mt-1">Hierarchical structure for data interchange</div>
          </div>
        </div>
        
        {#if downloads.xml}
          <div class="flex flex-wrap gap-3" role="group" aria-label="XML export actions">
            <SolidButton 
              color="green"
              href={downloads.xml.url}
              options={{ 
                download: downloads.xml.filename,
                'aria-label': `Download XML file: ${downloads.xml.filename}`
              }}
            >
              <Download class="w-4 h-4 mr-2" aria-hidden="true" />
              Download {downloads.xml.filename}
            </SolidButton>
            
            <SolidButton
              color="blue"
              onClick={() => copyToClipboard(downloads.xml!.data, 'xml')}
              options={{ 'aria-label': 'Copy XML data to clipboard' }}
            >
              {#if copiedStates.xml}
                <Check class="w-4 h-4 mr-2" aria-hidden="true" />
                Copied!
              {:else}
                <Copy class="w-4 h-4 mr-2" aria-hidden="true" />
                Copy Data
              {/if}
            </SolidButton>
          </div>
        {/if}
      </div>

      <!-- Instructions -->
      <div class="bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-700 rounded-lg p-6">
        <h3 class="font-medium text-blue-900 dark:text-blue-100 mb-3 flex items-center space-x-2">
          <span>How to use:</span>
        </h3>
        <ul class="text-blue-800 dark:text-blue-200 space-y-2" role="list">
          <li class="flex items-start space-x-2">
            <span class="font-medium text-green-600 dark:text-green-400">Download:</span>
            <span>Click the download button to save the file to your computer</span>
          </li>
          <li class="flex items-start space-x-2">
            <span class="font-medium text-blue-600 dark:text-blue-400">Copy:</span>
            <span>Click copy to put the data on your clipboard for pasting elsewhere</span>
          </li>
          <li class="flex items-start space-x-2">
            <span class="font-medium text-purple-600 dark:text-purple-400">Right-click:</span>
            <span>You can also right-click download links and choose "Save link as..."</span>
          </li>
        </ul>
      </div>
    </div>

    <div class="text-right mt-8">
      <SolidButton 
        color="gray" 
        onClick={closeModal}
      >
        Close
      </SolidButton>
    </div>
  </div>
</Modal>