import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from './types';

const { SecretFormField, FormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onPathChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      path: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  onAPIKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        apiKey: event.target.value,
      },
    });
  };

  onResetAPIKey = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        apiKey: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        apiKey: '',
      },
    });
  };

  render() {
    const { options } = this.props;
    const { jsonData, secureJsonFields } = options;
    const secureJsonData = (options.secureJsonData || {}) as MySecureJsonData;

    console.log('options', options);

    return (
        <div className="gf-form-group">
          <div className="gf-form">
            <FormField
                label="Project ID"
                labelWidth={6}
                inputWidth={20}
                onChange={this.onPathChange}
                value={jsonData.projectId || ''}
                placeholder="json field returned to frontend"
            />
          </div>

          <div className="gf-form-inline">
            <div className="gf-form">
              <SecretFormField
                  isConfigured={(secureJsonFields && secureJsonFields.apiKey) as boolean}
                  value={secureJsonData.publicKey || ''}
                  label="Public Key"
                  placeholder="secure json field (backend only)"
                  labelWidth={6}
                  inputWidth={20}
                  onReset={this.onResetAPIKey}
                  onChange={this.onAPIKeyChange}
              />
            </div>
          </div>

          <div className="gf-form-inline">
            <div className="gf-form">
              <SecretFormField
                  isConfigured={(secureJsonFields && secureJsonFields.apiKey) as boolean}
                  value={secureJsonData.privateKey || ''}
                  label="Private Key"
                  placeholder="secure json field (backend only)"
                  labelWidth={6}
                  inputWidth={20}
                  onReset={this.onResetAPIKey}
                  onChange={this.onAPIKeyChange}
              />
            </div>
          </div>
        </div>
    );
  }
}
