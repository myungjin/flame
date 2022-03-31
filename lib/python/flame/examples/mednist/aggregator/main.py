# Copyright 2022 Cisco Systems, Inc. and its affiliates
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# SPDX-License-Identifier: Apache-2.0


from ....common.util import install_packages

install_packages(['monai', 'sklearn', 'tqdm'])

# example cmd: python3 -m flame.examples.mednist.aggregator.main --rounds 3
# run the above command in flame/lib/python folder
if __name__ == "__main__":
    import argparse

    from .role import Aggregator

    parser = argparse.ArgumentParser(description='')
    parser.add_argument(
        '--rounds', type=int, default=1, help='number of training rounds'
    )

    args = parser.parse_args()

    aggregator = Aggregator(
        'flame/examples/mednist/aggregator/config.json',
        args.rounds,
    )
    aggregator.run()
