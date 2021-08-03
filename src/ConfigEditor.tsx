import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from './types';

const { SecretFormField, FormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onProjectIDChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      projectId: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  onPublicKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        publicKey: event.target.value,
      },
    });
  };

  onPrivateKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        privateKey: event.target.value,
      },
    });
  };

  onResetPublicKey = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        publicKey: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        publicKey: '',
      },
    });
  };

  onResetPrivateAPIKey = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        privateKey: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        privateKey: '',
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
            label="Project Id"
            labelWidth={6}
            inputWidth={20}
            onChange={this.onProjectIDChange}
            value={jsonData.projectId || ''}
            placeholder="Optional Project Id"
          />
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <SecretFormField
              isConfigured={(secureJsonFields && secureJsonFields.apiKey) as boolean}
              value={secureJsonData.publicKey || ''}
              label="Public Key"
              placeholder="Required UCloud Public Key"
              labelWidth={6}
              inputWidth={20}
              onReset={this.onResetPublicKey}
              onChange={this.onPublicKeyChange}
            />
          </div>
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <SecretFormField
              isConfigured={(secureJsonFields && secureJsonFields.apiKey) as boolean}
              value={secureJsonData.privateKey || ''}
              label="Private Key"
              placeholder="Required UCloud Private Key"
              labelWidth={6}
              inputWidth={20}
              onReset={this.onResetPrivateAPIKey}
              onChange={this.onPrivateKeyChange}
            />
          </div>
        </div>
      </div>
    );
  }
}
