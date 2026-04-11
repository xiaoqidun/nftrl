// Copyright 2026 肖其顿 (XIAO QI DUN)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

'use strict';
'require view';
'require form';
'require uci';
'require network';

return view.extend({
	load: function () {
		return Promise.all([
			uci.load('nftrl'),
			network.getHostHints()
		]);
	},

	render: function (data) {
		var hosts = data[1];
		var m, s, o;

		var enabled = uci.get('nftrl', 'global', 'enabled') === '1';
		var enableCb = E('input', {
			'type': 'checkbox',
			'checked': enabled ? '' : null,
			'click': function (ev) {
				uci.set('nftrl', 'global', 'enabled', ev.target.checked ? '1' : '0');
			}
		});

		var banner = E('div', {}, [
			E('h2', {}, E('a', { 'href': 'https://github.com/xiaoqidun/nftrl', 'target': '_blank' }, _('NFTRL'))),
			E('table', { 'class': 'table cbi-section-table' }, [
				E('tr', { 'class': 'tr cbi-section-table-row' }, [
					E('td', { 'class': 'td', 'width': '5%' }, enableCb),
					E('td', { 'class': 'td' },
						_('基于 MAC 地址的网络限速工具，规则运行在内核转发路径，仅影响外网流量，不限制局域网通信。'))
				])
			])
		]);

		m = new form.Map('nftrl');

		s = m.section(form.GridSection, 'device', _('限速设备列表'));
		s.anonymous = true;
		s.addremove = true;
		s.sortable = true;
		s.nodescriptions = true;

		o = s.option(form.Flag, 'enabled', _('启用'));
		o.default = '1';
		o.editable = true;
		o.width = '5%';
		o.rmempty = false;

		o = s.option(form.Value, 'comment', _('备注'));
		o.width = '25%';
		o.rmempty = true;

		o = s.option(form.Value, 'mac', _('MAC 地址'));
		o.datatype = 'macaddr';
		o.rmempty = false;
		o.width = '25%';

		if (hosts) {
			hosts.getMACHints().forEach(function (entry) {
				o.value(entry[0], entry[0] + (entry[1] ? ' (' + entry[1] + ')' : ''));
			});
		}

		o.validate = function (section_id, value) {
			if (!value)
				return _('MAC 地址不能为空');
			var mac = value.toLowerCase();
			var sections = uci.sections('nftrl', 'device');
			for (var i = 0; i < sections.length; i++) {
				if (sections[i]['.name'] === section_id)
					continue;
				if ((sections[i].mac || '').toLowerCase() === mac)
					return _('MAC 地址重复: ') + value;
			}
			return true;
		};

		o = s.option(form.Value, 'egress_limit', _('上传限速 (Kbps)'));
		o.datatype = 'uinteger';
		o.placeholder = _('不限制');
		o.width = '12.5%';
		o.rmempty = true;

		o = s.option(form.Value, 'ingress_limit', _('下载限速 (Kbps)'));
		o.datatype = 'uinteger';
		o.placeholder = _('不限制');
		o.width = '12.5%';
		o.rmempty = true;

		return m.render().then(function (mapNode) {
			return E('div', {}, [banner, mapNode]);
		});
	}
});